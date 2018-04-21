package main

import (
	"fmt"
	"io"
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
	filename := "/Users/timpringle/Desktop/videoplayback.mp4"
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
	joburi := "https://s3-eu-west-1.amazonaws.com/tim-training-thing/videoplayback.mp4"
	//mediaformat := "mp4"
	//languagecode := "en-US"

	var StrucMedia transcribeservice.Media
	StrucMedia.MediaFileUri = &joburi

	transcriber := transcribeservice.New(sess)
	/*
		transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
			TranscriptionJobName: &jobname,
			Media:                &StrucMedia,
			MediaFormat:          &mediaformat,
			LanguageCode:         &languagecode,
		})
	*/
	//Job processing will run async, so it's up to you how you deal with this.
	//For this one we'll take ten second naps in between checks of the status
	/*
		running := true
		for running == true {
			transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
				TranscriptionJobName: &jobname,
			})

			if err != nil {
				fmt.Println(err)
			}
			strStatus := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
			strJobname := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobName)

			fmt.Printf("Job %s is currently %s\n", strStatus, strJobname)
			if strStatus != "IN_PROGRESS" {
				running = false
			} else {
				time.Sleep(5 * time.Second)
			}
		}
	*/
	jobname = "01524-29669-18067-26570"
	transcriptionjoboutput, err := transcriber.GetTranscriptionJob(&transcribeservice.GetTranscriptionJobInput{
		TranscriptionJobName: &jobname,
	})

	if err != nil {
		fmt.Println(err)
	}

	strStatus := *(transcriptionjoboutput.TranscriptionJob.TranscriptionJobStatus)
	if strStatus == "COMPLETED" {
		uri := transcriptionjoboutput.TranscriptionJob.Transcript.TranscriptFileUri
		filename := "/Users/timpringle/Desktop/videoplayback.srt"
		DownloadFile(filename, *uri)
	}
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
