package main

import (
	"fmt"
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
	job_name := guid()
	job_uri := "https://s3-eu-west-1.amazonaws.com/tim-training-thing/videoplayback.mp4"
	media_format := "mp4"
	language_code := "en-US"

	var StrucMedia transcribeservice.Media
	StrucMedia.MediaFileUri = &job_uri

	transcriber := transcribeservice.New(sess)
	startjoboutput, err := transcriber.StartTranscriptionJob(&transcribeservice.StartTranscriptionJobInput{
		TranscriptionJobName: &job_name,
		Media:                &StrucMedia,
		MediaFormat:          &media_format,
		LanguageCode:         &language_code,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(startjoboutput)
	/*
		$null = Start-TRSTranscriptionJob -Media_MediaFileUri $s3uri -TranscriptionJobName $jobname -MediaFormat mp4 -LanguageCode en-US @AWSDefaultParameters

			#Job processing will run async, so it's up to you how you deal with this.
			#For this one we'll take ten second naps in between checks of the status
			$results = Get-TRSTranscriptionJob -TranscriptionJobName $jobname @AWSDefaultParameters

			While ($results.TranscriptionJobStatus -eq 'IN_PROGRESS') {
				Start-Sleep -Seconds 5
				$results = Get-TRSTranscriptionJob -TranscriptionJobName $jobname @AWSDefaultParameters
			}

			If ($results.TranscriptionJobStatus -eq 'COMPLETED') {
				$transcripturi = $results.Transcript.TranscriptFileUri
				Invoke-Webrequest -Uri $transcripturi -OutFile $resultsfile
				$output = Get-Content $resultsfile

				#Let's clear up the json file that was created
				Remove-Item -Path $resultsfile -Force

				#Output the results
				$output
	*/
}
