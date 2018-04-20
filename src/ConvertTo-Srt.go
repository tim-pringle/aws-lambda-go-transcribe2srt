package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

type Awstranscript struct {
	JobName   string `json:"jobName"`
	Accountid string `json:"accountId"`
	Results   Result `json:"results"`
	Status    string `json:"status"`
}

type Result struct {
	Transcripts []Transcript `json:"transcripts"`
	Items       []Item       `json:"items"`
}

type Transcript struct {
	Transcript string `json:"transcript"`
}

type Item struct {
	Starttime      string        `json:"start_time"`
	Endtime        string        `json:"end_time"`
	Alternatives   []Alternative `json:"alternatives"`
	Classification string        `json:"type"`
}

type Alternative struct {
	Confidence string `json:"confidence"`
	Content    string `json:"content"`
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("/Users/timpringle/Downloads/asrOutput.json")

	//If there's an error, print the error
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened asrOutput.json")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read the opened file as a byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// initialize our variable to hold the json
	var awstranscript Awstranscript

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'awstranscript' which we defined above
	json.Unmarshal(byteValue, &awstranscript)

	// some sample examples of outputing the json
	fmt.Println("Job Name: " + awstranscript.JobName)
	fmt.Println("Transcript: " + awstranscript.Results.Transcripts[0].Transcript)
	fmt.Println("Item: " + string(awstranscript.Results.Items[0].Starttime))
	fmt.Println("Content: " + awstranscript.Results.Items[0].Alternatives[0].Content)

	var transcription []Item
	transcription = awstranscript.Results.Items

	var index, sequence int = 0, 0
	var srtinfo, subdetail, subtitle, sttime, classification, text string
	var strlen int
	subtitle += "!"

	for index = 0; index < len(transcription); srtinfo += subdetail {
		strlen = 0

		//Grab the start time of the item
		sttime = transcription[index].Starttime
		fsttime, err := strconv.ParseFloat(sttime, 64)

		sttime = timestr(fsttime)

		fmt.Println(sttime)
		if err != nil {
			fmt.Println(err)
		}

		//sttime := [timespan]::FromSeconds($sttime)
		//starttime = "{0:hh}:{0:mm}:{0:ss},{0:fff}" -f $sttime
		//subtitle := ""
		sequence++
		subtitle += "!"
		subdetail += "!"
		srtinfo += "!"
		classification += "!"
		text += "!"
		strlen++
		//firstrow := true

		index++

	}
}

func timestr(numerator float64) (timestring string) {

	var h = 3600
	var m = 60
	var s = 1

	integer, frac := math.Modf(numerator)

	integerpart := int(integer)

	hours := integerpart / h
	remainder := integerpart % h

	minutes := remainder / m
	remainder = remainder % m

	seconds := remainder / s
	stringfrac := strconv.FormatFloat(frac, 'f', 3, 64)
	runes := []rune(stringfrac)
	safeSubstring := string(runes[1:len(stringfrac)])

	timestring = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	timestring += safeSubstring
	return
}
