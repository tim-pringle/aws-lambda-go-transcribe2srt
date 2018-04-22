package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

//GUID - generates a unique identifier
func GUID() (guid string) {
	ad, err := time.Parse("02-01-2006", "01-01-1970")

	if err != nil {
	}

	timesince := time.Since(ad).Nanoseconds()
	strsince := strconv.FormatInt(timesince, 10)
	guid = fmt.Sprintf("0" + strsince[0:4] + "-" + strsince[4:9] + "-" + strsince[9:14] + "-" + strsince[14:19])
	return
}

//DownloadFile -  Downloads a file from the web
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

// Generate an SRT format time string
func getsrttime(numerator float64) (timestring string) {

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
