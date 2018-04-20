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

	var transcription []Item
	transcription = awstranscript.Results.Items

	var index, sequence int = 0, 0
	var srtinfo, subdetail, subtitle, sttime, classification, text, entime string
	var strlen int
	var firstrow bool

	for index = 0; index < len(transcription); {
		//Variable initiation for length of subtitle text, sequence number if its the first line and the subtitle text

		sequence++
		firstrow = true
		subtitle = ""

		//Grab the start time, convert it to a number, then convert the number an SRT valid time string
		sttime = transcription[index].Starttime
		fsttime, err := strconv.ParseFloat(sttime, 64)
		if err != nil {
			fmt.Println(err)
		}
		sttime = timestr(fsttime)

		/*Repeat this until we have either reached the last item in results
		#or the length of the lines we are reading is greater than 64 characters */

		for strlen = 0; (strlen < 64) && (index < len(transcription)); {
			text = transcription[index].Alternatives[0].Content
			strlen += len(text)

			switch classification {

			case "punctuation":
				if len(subtitle) > 0 {
					runes := []rune(subtitle)
					subtitle = string(runes[1 : len(subtitle)-1])
				} else {
					subtitle += text
				}
			default:
				subtitle += (text + " ")
			}

			//If the length of the current string is greater than 32 and this
			//is the first line of the sequence, then add a return character to it

			if (strlen > 32) && (firstrow == true) {
				subtitle += "\n"
				firstrow = false
			}

			/*If the last character is a '.', then we need to set
			the end time attribute to the previous indexes one
			since punctuation characters to not have a time stamp*/

			if classification == "punctuation" {
				entime = transcription[index-1].Endtime
			} else {
				entime = transcription[index].Endtime
			}

			fsttime, err = strconv.ParseFloat(entime, 64)
			entime = timestr(fsttime)

			index++
		}
		//Setup the string that is refers to these two
		//lines in SRT format

		subdetail = fmt.Sprintf("\n%d\n%s --> %s\n%s\n", sequence, sttime, entime, subtitle)

		//Append this to the existing string
		srtinfo += subdetail

	}

	//#Now output the results to our .srt file
	//$srtinfo | Set-Content $DestinationPath -Force
	fmt.Println(srtinfo)
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
