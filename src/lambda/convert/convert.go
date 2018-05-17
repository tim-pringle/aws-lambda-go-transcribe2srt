package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/tim-pringle/go-aws/transcribe"

	"github.com/tim-pringle/go-misc/misc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

var (
	ErrJobnameNotProvided = errors.New("No job number has been provided HTTP body")
	ErrTranscribeRunning  = errors.New("The transcription job is running")
	ErrTranscribeFailure  = errors.New("The transcription job has failed")
	ErrJobDoesntExist     = errors.New("The transcripting job does not exist")
)

//Convert is the script that is called by the appropriate requestors to provide an srt file
func Convert(jobname string) (string, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing conversion request")

	sess, _ := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-west-1")},
		Profile: "development",
	})

	log.Printf("Job name : %s", jobname)
	log.Printf("Creating new session")
	transcriber := transcribeservice.New(sess)

	log.Printf("Getting transcription job")
	transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
		TranscriptionJobName: &jobname,
	})

	if err != nil {
		ErrMsg := errors.New(err.Error())
		log.Printf("%s", err.Error())
		return "", ErrMsg
	}

	strStatus := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
	log.Printf("Job status is %s", strStatus)

	if strStatus == "FAILED" {
		return "", ErrTranscribeFailure
	}
	if strStatus == "IN_PROGRESS" {
		return "", ErrTranscribeRunning
	}

	var uri *string
	uri = transcriptionjoboutput.TranscriptionJob.Transcript.TranscriptFileUri
	log.Printf("URI is %s", *uri)

	response, _ := http.Get(*uri)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	//If there's an error, print the error
	if err != nil {
		fmt.Println(err)
	}

	// initialize our variable to hold the json
	var awstranscript transcribe.Awstranscript

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'awstranscript' which we defined above
	json.Unmarshal(body, &awstranscript)

	var transcription []transcribe.Item
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
		sttime = misc.Getsrttime(fsttime)

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
			entime = misc.Getsrttime(fsttime)

			index++
		}
		//Setup the string that is refers to these two
		//lines in SRT format

		subdetail = fmt.Sprintf("\n%d\n%s --> %s\n%s\n", sequence, sttime, entime, subtitle)

		//Append this to the existing string
		srtinfo += subdetail

	}

	log.Printf(srtinfo)

	return srtinfo, nil
}
