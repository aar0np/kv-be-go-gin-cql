package controllers

import (
	"errors"
	"killrvideo/go-backend-astra-cql/models"
	repo "killrvideo/go-backend-astra-cql/repository"
	"net/http"
	"os"
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	Access   string
	Refresh  string
	JTIAcc   string
	JTIRef   string
	ExpAcc   time.Time
	ExpRef   time.Time
	UserID   string
	Issuer   string
	Audience string
}

type AuthController struct {
	authDAL repo.AuthDAL
}

func NewAuthController(session *apachegocql.Session) *AuthController {
	return &AuthController{
		authDAL: *repo.NewAuthDAL(session),
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var newUserReg models.UserRegistrationRequest

	if err1 := c.BindJSON(&newUserReg); err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	var creds models.UserCredentials
	var user models.User
	var jwtResp models.JwtResponse

	// generate credentials
	var newUserId = uuid.New()
	creds.Email = newUserReg.Email
	creds.Userid = apachegocql.UUID(newUserId)
	creds.Password = hashPassword(newUserReg.Password)
	creds.AccountLocked = false

	// generate user
	user.Email = newUserReg.Email
	user.Userid = apachegocql.UUID(newUserId)
	user.FirstName = newUserReg.FirstName
	user.LastName = newUserReg.LastName
	user.CreatedDate = time.Now()
	user.LastLoginDate = time.Now()

	// save to DB
	ac.authDAL.SaveUser(user)
	ac.authDAL.SaveUserCreds(creds)

	// gen token
	token, err3 := issueToken(newUserId.String(), newUserReg.Email)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	jwtResp.Email = newUserReg.Email
	jwtResp.UserID = newUserId.String()
	jwtResp.Token = token.Access

	c.JSON(http.StatusOK, jwtResp)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err1 := c.BindJSON(&req); err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	user, err2 := ac.authDAL.GetUserCredsByEmail(req.Email)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	hashedPassword := user.Password
	id := user.Userid

	if !validatePassword(req.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password."})
		return
	}

	token, err3 := issueToken(id.String(), req.Email)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3.Error()})
		return
	}

	var jwtResp models.JwtResponse
	jwtResp.Email = req.Email
	jwtResp.UserID = id.String()
	jwtResp.Token = token.Access

	c.JSON(http.StatusOK, jwtResp)
}

func (ac *AuthController) GetUser(c *gin.Context) {
	id, err1 := apachegocql.ParseUUID(c.Param("id"))

	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	user, err2 := ac.authDAL.GetUserById(id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func (ac *AuthController) GetCurrentUser(c *gin.Context) {
	// parse UserID from request
	userid, err1 := getUserIdFromAuth(c)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	// get User from DB
	user, err2 := ac.authDAL.GetUserById(userid)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func getUserIdFromAuth(c *gin.Context) (apachegocql.UUID, error) {
	token := c.Query("token")

	claims, err1 := parseWithSecret(token)
	if err1 != nil {
		return apachegocql.MustRandomUUID(), err1
	}

	// get UserID from Subject
	uuid, err2 := apachegocql.ParseUUID(claims.Subject)
	return uuid, err2
}

func validatePassword(password string, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == nil {
		return true
	} else {
		return false
	}
}

func hashPassword(password string) string {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(hashed)
}

func issueToken(userID string, email string) (*Tokens, error) {
	now := time.Now().UTC()

	key := os.Getenv("JWT_KEY")

	t := &Tokens{
		UserID:   userID,
		JTIAcc:   uuid.NewString(),
		JTIRef:   uuid.NewString(),
		ExpAcc:   now.Add(15 * time.Minute),
		ExpRef:   now.Add(7 * 24 * time.Hour),
		Issuer:   "kv-be-go-gin-cql",
		Audience: "killrvideo-react-frontend",
	}

	acc := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIAcc,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{userID, email},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpAcc),
	})

	ref := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ID:        t.JTIRef,
		Issuer:    t.Issuer,
		Audience:  jwt.ClaimStrings{userID, email},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(t.ExpRef),
	})

	var err error
	t.Access, err = acc.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}
	t.Refresh, err = ref.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func parseWithSecret(tokenStr string) (*jwt.RegisteredClaims, error) {

	secret := os.Getenv("JWT_KEY")

	if secret == "" {
		return nil, errors.New("jwt secret not configured")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Extra safety: ensure HMAC family
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
