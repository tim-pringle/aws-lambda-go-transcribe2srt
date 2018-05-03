# aws-lambda-go-transcribe2srt
This is a work-in-progress project, which I've been working on to learn Golang. It uses several AWS services to convert the dialog in media file to SRT format

The entire project is written in GO, and makes use of the following AWS resources:

* S3
* Lambda
* Transcribe Service

## Transcribing
This uses the *transcribe.go* script. 
Only S3 is supported for the source location of the media files. There's no upload functionality in the script.

Make a Lambda function for the Transcribe script. Link this to an S3 bucket. Create an Event watch to kick off the function when a CreateObject event has occurred. This handles the Transcribe job itself.

## Obtaining the SRT
This uses the *convert.go* script

Two sources are, in varying levels of completeness, in place for triggering the job to convert and provide an SRT file from a completed Transcribe job:

* API Gateway
* Lex

Also in place is processing of an Alexa skill. This only reads out the transcript from the job, since it doesn't really make sense to read it out in SRT format.

## Known Issues
The structs, method for handling data, and general layout of the code suck big time...

The end time entry for the last sequence in the SRT file is always wrong.

## Contributing
Bug reports and contributions for the project are always welcome at the repo on [GitHub](https://github.com/tim-pringle/aws-lambda-go-transcribe2srt).

##History
22/04/2018 : First operating version
24/04/2018 : API Gateway and Lex support
03/05/2018 : Alexa integration added

Developed on OSX 10.13.3 with Go v1.10

##License
Apache 2.0 (see LICENSE)