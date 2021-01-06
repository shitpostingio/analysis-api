package client

import (
	"bytes"
	"encoding/json"
	"github.com/shitpostingio/analysis-commons/structs"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

// PerformFingerprintRequest performs a request to the fingerprinting service.
func PerformFingerprintRequest(file io.Reader, fileName, endpoint, authorization string) (data structs.FingerprintResponse, errorString string) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		errorString = err.Error()
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		errorString = err.Error()
		return
	}

	err = writer.Close()
	if err != nil {
		errorString = err.Error()
		return
	}

	request, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		errorString = err.Error()
		return
	}

	// We want to send data with the multipart/form-data Content-Type
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-shitposting-key", authorization)
	client := http.Client{Timeout: time.Second * 30}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Request failed: ", err)
		errorString = err.Error()
		return
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Println("Analysis.PerformFingerprintRequest: unable to close response body", err)
		}
	}()

	bodyResult, err := ioutil.ReadAll(response.Body)
	log.Println("Body response: ", string(bodyResult))
	if err != nil {
		errorString = err.Error()
		return
	}

	var ar structs.FingerprintResponse
	err = json.Unmarshal(bodyResult, &ar)
	if err != nil {
		errorString = err.Error()
		log.Println("PerformRequest: error while unmarshaling ", err)
		return
	}

	return ar, errorString

}
