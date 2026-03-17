package v1

// @Summary     Translate text
// @Description Translate text from source language to destination and store in history
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body request.Translation true "Translation request"
// @Success     200 {object} response.Translation
// @Failure     400 {object} response.Error
// @Router      /v1/translation/do [post]
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
