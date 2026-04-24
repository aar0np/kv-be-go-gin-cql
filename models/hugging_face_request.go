package models

type HuggingFaceRequest struct {
	Model string `json:"model"`
	Text  string `json:"text"`
}
