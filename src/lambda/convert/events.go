package main

type AlexaRequest struct {
	Version string `json:"version"`
	Session struct {
		New         bool   `json:"new"`
		SessionID   string `json:"sessionId"`
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		Attributes struct {
			Key string `json:"key"`
		} `json:"attributes"`
		User struct {
			UserID      string `json:"userId"`
			AccessToken string `json:"accessToken"`
			Permissions struct {
				ConsentToken string `json:"consentToken"`
			} `json:"permissions"`
		} `json:"user"`
	} `json:"session"`
	Context struct {
		System struct {
			Device struct {
				DeviceID            string `json:"deviceId"`
				SupportedInterfaces struct {
					AudioPlayer struct {
					} `json:"AudioPlayer"`
				} `json:"supportedInterfaces"`
			} `json:"device"`
			Application struct {
				ApplicationID string `json:"applicationId"`
			} `json:"application"`
			User struct {
				UserID      string `json:"userId"`
				AccessToken string `json:"accessToken"`
				Permissions struct {
					ConsentToken string `json:"consentToken"`
				} `json:"permissions"`
			} `json:"user"`
			APIEndpoint    string `json:"apiEndpoint"`
			APIAccessToken string `json:"apiAccessToken"`
		} `json:"System"`
		AudioPlayer struct {
			PlayerActivity       string `json:"playerActivity"`
			Token                string `json:"token"`
			OffsetInMilliseconds int    `json:"offsetInMilliseconds"`
		} `json:"AudioPlayer"`
	} `json:"context"`
	Request struct {
	} `json:"request"`
}

type AlexaResponse struct {
	Version           string `json:"version"`
	SessionAttributes struct {
		Key string `json:"key"`
	} `json:"sessionAttributes"`
	Response struct {
		OutputSpeech struct {
			Type string `json:"type"`
			Text string `json:"text"`
			Ssml string `json:"ssml"`
		} `json:"outputSpeech"`
		Card struct {
			Type    string `json:"type"`
			Title   string `json:"title"`
			Content string `json:"content"`
			Text    string `json:"text"`
			Image   struct {
				SmallImageURL string `json:"smallImageUrl"`
				LargeImageURL string `json:"largeImageUrl"`
			} `json:"image"`
		} `json:"card"`
		Reprompt struct {
			OutputSpeech struct {
				Type string `json:"type"`
				Text string `json:"text"`
				Ssml string `json:"ssml"`
			} `json:"outputSpeech"`
		} `json:"reprompt"`
		Directives []struct {
			Type string `json:"type"`
		} `json:"directives"`
		ShouldEndSession bool `json:"shouldEndSession"`
	} `json:"response"`
}

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
