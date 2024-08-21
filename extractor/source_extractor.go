package extractor

import (
	"errors"
	"io"
	"log"
	"net/http"
)

func Extract(url string) (string, error) {

	// Send a GET request
	var resp, getErr = http.Get(url)
	if getErr != nil {
		return "", getErr
	}

	// Close the response body on function exit
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch URL: status code %d\n", resp.StatusCode)
		return "", errors.New("failed to fetch URL")
	}

	// Read the response body
	var body, readErr = io.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("Failed to read response body: %s\n", readErr)
		return "", readErr
	}

	// Return the response body as a string
	return string(body), nil
}
