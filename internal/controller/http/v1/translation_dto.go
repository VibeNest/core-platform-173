package v1

// RequestDTO для метода перевода
type doTranslateRequest struct {
	Source      string `json:"source" validate:"required" example:"en"`
	Destination string `json:"destination" validate:"required" example:"ru"`
	Text        string `json:"text" validate:"required" example:"Hello world"`
}

// ResponseDTO для элемента истории
type historyResponse struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Original    string `json:"original"`
	Translation string `json:"translation"`
}
