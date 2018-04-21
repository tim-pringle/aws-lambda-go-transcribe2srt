package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
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

func guid() (guid string) {
	ad, err := time.Parse("02-01-2006", "01-01-1970")

	if err != nil {
	}

	timesince := time.Since(ad).Nanoseconds()
	strsince := strconv.FormatInt(timesince, 10)
	guid = fmt.Sprintf("0" + strsince[0:4] + "-" + strsince[4:9] + "-" + strsince[9:14] + "-" + strsince[14:19])
	return
}

func main() {
	//Set the S3 uri prefix and create a unique guid to be used as the job name

	//resultsfile := "/Users/timpringle/Downloads/asrOutput.json"
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

	//Let's transcribe!!
	jobname := guid()
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
	//For this one we'll take ten second naps in between checks of the status

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

	out, err := os.Create("/Users/timpringle/Downloads/AWS Summit San Francisco 2018 - Amazon Transcribe Now Generally Available.srt")
	defer out.Close()
	_, err = out.WriteString(srtinfo)

	fmt.Println(srtinfo)
}

func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err

	}

	return nil
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
