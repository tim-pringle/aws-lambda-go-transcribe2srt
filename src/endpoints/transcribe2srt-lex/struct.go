package main

//LexResponse is use to provide a valid response to a request from a Lex endpoint
type LexResponse struct {
	SessionAttributes struct {
	} `json:"sessionAttributes"`
	DialogAction struct {
		Type             string `json:"type"`
		FulfillmentState string `json:"fulfillmentState"`
		Message          struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"message"`
	} `json:"dialogAction"`
}
