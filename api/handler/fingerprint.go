package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shitpostingio/analysis-api/api/cache"
	"github.com/shitpostingio/analysis-commons/downloader"
	"github.com/shitpostingio/analysis-commons/structs"
	"github.com/shitpostingio/fingerprint-microservice/client"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"io"
	"net/http"
)

var fingerprintGroup singleflight.Group

// HandleFingerprint handles a fingerprint request.
// It will check the redis cache first, if the element is not found
// it will then proceed to download the file and fingerprint it.
func (h *Handler) HandleFingerprint(w http.ResponseWriter, r *http.Request, d downloader.Downloader) {

	// Check cache for a hit first.
	rVars := mux.Vars(r)
	id := rVars["id"]

	data, err := cache.GetFingerprint(id)
	if err == nil {
		log.Debugln("Cache hit for fingerprint ", id)
		encodeFingerprintResponse(w, data)
		return
	}

	// Perform this part in a shared manner
	result, err, _ := fingerprintGroup.Do(id, func() (i interface{}, err error) {

		mediaType := rVars["type"]
		filename, reader, err := d.Download(id, mediaType, r)
		if err != nil {
			log.Println("requestFingerprint: error while downloading file: ", err)
			return
		}

		defer func() {
			err = reader.Close()
			if err != nil {
				log.Println("requestFingerprint: error while trying to close download reader:", err)
			}
		}()

		log.Debugln("Performing fingerprint request ", id)
		return h.performFingerprint(reader, id, filename, mediaType)

	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		encodeFingerprintResponse(w, result.(structs.FingerprintResponse))
	}

}

func (h *Handler) performFingerprint(reader io.Reader, id, filename, mediaType string) (f structs.FingerprintResponse, err error) {

	endpoint := fmt.Sprintf("%s/%s", h.FingerprintEndpoint, mediaType)
	response, errString := client.PerformRequest(reader, filename, endpoint)
	if errString != "" {
		log.Println("performFingerprint: error while performing request to fingerprint service: ", errString)
		return response, errors.New(errString)
	}

	cacheErr := cache.PutFingerprint(id, response)
	if cacheErr != nil {
		log.Println("performFingerprint: unable to save fingerprint response to cache: ", err)
	}

	return response, nil

}

func encodeFingerprintResponse(w http.ResponseWriter, f structs.FingerprintResponse) {

	err := json.NewEncoder(w).Encode(f)
	if err != nil {
		log.Println("encodeFingerprintResponse: unable to send response:", err)
	}

}
