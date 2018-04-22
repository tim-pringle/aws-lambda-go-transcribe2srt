package main

//Awstranscript - Top level struc for an AWS transcript job output
type Awstranscript struct {
	JobName   string `json:"jobName"`
	Accountid string `json:"accountId"`
	Results   Result `json:"results"`
	Status    string `json:"status"`
}

//Result - Result structure
type Result struct {
	Transcripts []Transcript `json:"transcripts"`
	Items       []Item       `json:"items"`
}

//Transcript - Transcription
type Transcript struct {
	Transcript string `json:"transcript"`
}

//Item - Individual translation word/punctuation record
type Item struct {
	Starttime      string        `json:"start_time"`
	Endtime        string        `json:"end_time"`
	Alternatives   []Alternative `json:"alternatives"`
	Classification string        `json:"type"`
}

// Alternative - Actual translated word and confidence of accuracy
type Alternative struct {
	Confidence string `json:"confidence"`
	Content    string `json:"content"`
}
