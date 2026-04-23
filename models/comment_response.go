package models

type CommentResponse struct {
	Data       []Comment  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func NewCommentResponse() *CommentResponse {
	return &CommentResponse{
		Data: make([]Comment, 0),
	}
}
