package models

type HuggingFaceResponse struct {
	// Embedding [384]float64 is specific to IBM Granite 30m English model
	Embedding       [384]float32 `json:"Embedding"`
	Dim             int          `json:"dim"`
	Model           string       `json:"model"`
	TrustRemoteCode bool         `json:"trust_remote_code"`
	Predefined      bool         `json:"predefined"`
}
