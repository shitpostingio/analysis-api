package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shitpostingio/analysis-api/api/cache"
	"github.com/shitpostingio/analysis-commons/client"
	"github.com/shitpostingio/analysis-commons/downloader"
	"github.com/shitpostingio/analysis-commons/structs"
	"golang.org/x/sync/singleflight"
	"io"
	"log"
	"net/http"
)

var nsfwGroup singleflight.Group

// HandleNSFW handles a NSFW request.
// It will check the redis cache first, if the element is not found
// it will then proceed to download the file and check it for NSFW.
func (h *Handler) HandleNSFW(w http.ResponseWriter, r *http.Request, d downloader.Downloader) {

	// Check cache for a hit first.
	rVars := mux.Vars(r)
	id := rVars["id"]
	data, err := cache.GetNSFW(id)
	if err == nil {
		encodeNSFWResponse(w, data)
		return
	}

	result, err, _ := nsfwGroup.Do(id, func() (i interface{}, err error) {

		// Download the file otherwise
		mediaType := rVars["type"]
		filename, reader, err := d.Download(id, mediaType, r)
		if err != nil {
			log.Println("HandleNSFW: error while downloading file:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			err = reader.Close()
			if err != nil {
				log.Println("HandleNSFW: handle while closing reader:", err)
			}
		}()

		return h.performNSFW(reader, id, filename, mediaType)

	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		encodeNSFWResponse(w, result.(structs.NSFWResponse))
	}

}

func (h *Handler) performNSFW(reader io.Reader, id, filename, mediaType string) (f structs.NSFWResponse, err error) {

	endpoint := fmt.Sprintf("%s/%s", h.NSFWEndpoint, mediaType)
	response, err := client.PerformRequest(reader, filename, endpoint, "", "")
	if err != nil {
		log.Println("performNSFW: error while performing request to NSFW service: ", err)
		return structs.NSFWResponse{}, err
	}

	cacheErr := cache.PutNSFW(id, response.NSFW)
	if cacheErr != nil {
		log.Println("performNSFW: unable to save NSFW response to cache: ", err)
	}

	return response.NSFW, nil

}

func encodeNSFWResponse(w http.ResponseWriter, f structs.NSFWResponse) {

	err := json.NewEncoder(w).Encode(f)
	if err != nil {
		log.Println("encodeNSFWResponse: unable to send response:", err)
	}

}
