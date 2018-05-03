package main

import "time"

type AlexaRequest struct {
	Context struct {
		AudioPlayer struct {
			PlayerActivity string `json:"playerActivity"`
		} `json:"AudioPlayer"`
		Display struct {
		} `json:"Display"`
		System struct {
			APIAccessToken string `json:"apiAccessToken"`
			APIEndpoint    string `json:"apiEndpoint"`
			Application    struct {
				ApplicationID string `json:"applicationId"`
			} `json:"application"`
			Device struct {
				DeviceID            string `json:"deviceId"`
				SupportedInterfaces struct {
					AudioPlayer struct {
					} `json:"AudioPlayer"`
					Display struct {
						MarkupVersion   string `json:"markupVersion"`
						TemplateVersion string `json:"templateVersion"`
					} `json:"Display"`
				} `json:"supportedInterfaces"`
			} `json:"device"`
			User struct {
				UserID string `json:"userId"`
			} `json:"user"`
		} `json:"System"`
	} `json:"context"`
	Request struct {
		Locale                     string    `json:"locale"`
		RequestID                  string    `json:"requestId"`
		ShouldLinkResultBeReturned bool      `json:"shouldLinkResultBeReturned"`
		Timestamp                  time.Time `json:"timestamp"`
		Type                       string    `json:"type"`
	} `json:"request"`
	Session struct {
		Application struct {
			ApplicationID string `json:"applicationId"`
		} `json:"application"`
		New       bool   `json:"new"`
		SessionID string `json:"sessionId"`
		User      struct {
			UserID string `json:"userId"`
		} `json:"user"`
	} `json:"session"`
	Version string `json:"version"`
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
