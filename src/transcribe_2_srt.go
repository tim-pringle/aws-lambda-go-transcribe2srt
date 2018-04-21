package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

func main() {
	bucket := "tim-training-thing"
	filename := "/Users/timpringle/Downloads/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.mp4"
	//Lets set default parameter values for the aws cmdlets

	sess, _ := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("eu-west-1")},
		Profile: "development",
	})

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filename)
	key := filepath.Base(file.Name())

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		fmt.Println(err)
	}

	//Set the S3 uri for the file location and create a unique guid to be used as the job name
	jobname := GUID()
	joburi := "https://s3-eu-west-1.amazonaws.com/tim-training-thing/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.mp4"
	mediaformat := "mp4"
	languagecode := "en-US"

	var StrucMedia transcribeservice.Media
	StrucMedia.MediaFileUri = &joburi

	transcriber := transcribeservice.New(sess)

	transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
		TranscriptionJobName: &jobname,
		Media:                &StrucMedia,
		MediaFormat:          &mediaformat,
		LanguageCode:         &languagecode,
	})

	//Job processing will run async, so it's up to you how you deal with this.
	//For this one we'll take 5 second naps in between checks of the status

	running := true
	var strStatus, strJobname string

	for running == true {
		transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
			TranscriptionJobName: &jobname,
		})

		if err != nil {
			fmt.Println(err)
		}
		strStatus = *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
		strJobname = *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobName)

		fmt.Printf("Job %s is currently %s\n", strJobname, strStatus)
		if strStatus != "IN_PROGRESS" {
			running = false
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	//The job is no longer running. Let's check the status
	transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
		TranscriptionJobName: &strJobname,
	})

	if err != nil {
		fmt.Println(err)
	}

	strStatus = *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
	if strStatus == "COMPLETED" {
		uri := transcriptionjoboutput.TranscriptionJob.Transcript.TranscriptFileUri
		filename := "/Users/timpringle/Downloads/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.json"
		DownloadFile(filename, *uri)
	}

	// Open our jsonFile
	jsonFile, err := os.Open("/Users/timpringle/Downloads/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.json")

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
		sttime = getsrttime(fsttime)

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
			entime = getsrttime(fsttime)

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

	out, err := os.Create("/Users/timpringle/Downloads/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.srt")
	defer out.Close()
	_, err = out.WriteString(srtinfo)

	fmt.Println(srtinfo)
}
