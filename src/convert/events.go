package main

import "time"

type AlexaRequest struct {
	Version string `json:"version"`
	Session struct {
		New         bool   `json:"new"`
		SessionID   string `json:"sessionId"`
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		User struct {
			UserID string `json:"userId"`
		} `json:"user"`
	} `json:"session"`
	Context struct {
		AudioPlayer struct {
			PlayerActivity string `json:"playerActivity"`
		} `json:"AudioPlayer"`
		Display struct {
			Token string `json:"token"`
		} `json:"Display"`
		System struct {
			Application struct {
				ApplicationID string `json:"applicationId"`
			} `json:"application"`
			User struct {
				UserID string `json:"userId"`
			} `json:"user"`
			Device struct {
				DeviceID            string `json:"deviceId"`
				SupportedInterfaces struct {
					AudioPlayer struct {
					} `json:"AudioPlayer"`
					Display struct {
						TemplateVersion string `json:"templateVersion"`
						MarkupVersion   string `json:"markupVersion"`
					} `json:"Display"`
				} `json:"supportedInterfaces"`
			} `json:"device"`
			APIEndpoint    string `json:"apiEndpoint"`
			APIAccessToken string `json:"apiAccessToken"`
		} `json:"System"`
	} `json:"context"`
	Request struct {
		Type      string    `json:"type"`
		RequestID string    `json:"requestId"`
		Timestamp time.Time `json:"timestamp"`
		Locale    string    `json:"locale"`
		Intent    struct {
			Name               string `json:"name"`
			ConfirmationStatus string `json:"confirmationStatus"`
			Slots              struct {
				Jobnumber struct {
					Name               string `json:"name"`
					Value              string `json:"value"`
					ConfirmationStatus string `json:"confirmationStatus"`
				} `json:"jobnumber"`
			} `json:"slots"`
		} `json:"intent"`
		DialogState string `json:"dialogState"`
	} `json:"request"`
}

type SessionAttributes struct {
	Key string `json:"key,omitempty"`
}

type Card struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Text    string `json:"text,omitempty"`
	Image   Image  `json:"image,omitempty"`
}

type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech,omitempty"`
}

type OutputSpeech struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	Ssml string `json:"ssml,omitempty"`
}

type Directives struct {
	Type          string `json:"type,omitempty"`
	SlotToElicit  string `json:"slotToElicit,omitempty"`
	UpdatedIntent string `json:"updatedIntent"`
}

type Response struct {
	OutputSpeech     *OutputSpeech     `json:"outputSpeech,omitempty"`
	Card             *Card             `json:"card,omitempty"`
	Reprompt         *Reprompt         `json:"reprompt,omitempty"`
	ShouldEndSession bool              `json:"shouldEndSession"`
	Directives       *[]DialogDelegate `json:"directives,omitempty"`
}

type AlexaResponse struct {
	Version           string             `json:"version,omitempty"`
	SessionAttributes *SessionAttributes `json:"sessionAttributes,omitempty"`
	Response          Response           `json:"response,omitempty"`
}

type DialogDelegate struct {
	Type          string `json:"type"`
	UpdatedIntent *struct {
		Name               string `json:"name"`
		ConfirmationStatus string `json:"confirmationStatus"`
		Slots              struct {
			String struct {
				Name               string `json:"name"`
				Value              string `json:"value"`
				ConfirmationStatus string `json:"confirmationStatus"`
			} `json:"string"`
		} `json:"slots"`
	}
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

type JobResponse struct {
	Jobnumber struct {
		Name               string `json:"name"`
		Value              string `json:"value"`
		ConfirmationStatus string `json:"confirmationStatus"`
	} `json:"jobnumber"`
}
