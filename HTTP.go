package valitor

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func SendRequest(envelope string, method string, url string) ([]byte, error) {

	client := &http.Client{}
	// build a new request, but not doing the POST yet
	log.Println("SENDING TO:", method, url)
	log.Println(envelope)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(envelope)))
	if err != nil {
		return nil, err
	}

	// you can then set the Header here
	// I think the content-type should be "application/xml" like json...
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	// now POST it
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// read the response body to a variable
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("body!")
	log.Println(resp.StatusCode)
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	//print raw response body for debugging purposes
	return bodyBytes, nil
}

func send(url, method, body string) (results []byte, err error) {
	results, requestError := SendRequest(body, "POST", url)
	if requestError != nil {
		return results, requestError
	}

	return results, nil
}

func SendJSON(data []byte, method string, url string) ([]byte, int, error) {

	client := &http.Client{}
	// build a new request, but not doing the POST yet
	log.Println("SENDING TO:", method, url)
	log.Println(string(data))
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, err
	}

	// you can then set the Header here
	// I think the content-type should be "application/xml" like json...
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	// now POST it
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	// read the response body to a variable
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	log.Println("body!")
	log.Println(resp.StatusCode)
	bodyString := string(bodyBytes)
	log.Println(bodyString)
	//print raw response body for debugging purposes
	return bodyBytes, resp.StatusCode, nil
}
