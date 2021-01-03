package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shitpostingio/analysis-api/api/cache"
	"github.com/shitpostingio/analysis-commons/downloader"
	"github.com/shitpostingio/analysis-commons/structs"
	"golang.org/x/sync/singleflight"
	"io/ioutil"
	"log"
	"net/http"
)

var analysisGroup singleflight.Group

// HandleAnalysis handles a analysis request.
// It will check the redis cache first, if the element is not found
// it will then proceed to download the file and analyze it concurrently.
func (h *Handler) HandleAnalysis(w http.ResponseWriter, r *http.Request, d downloader.Downloader) {

	// Check cache for hits.
	rVars := mux.Vars(r)
	id := rVars["id"]

	gotF, gotN, response := checkCacheForAnalysis(id)
	if gotF && gotN {
		encodeAnalysisResponse(w, response)
		return
	}

	result, err, _ := analysisGroup.Do(id, func() (i interface{}, err error) {
		mediaType := rVars["type"]
		return h.requestAnalysis(r, d, id, mediaType, gotF, gotN)
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		encodeAnalysisResponse(w, result.(structs.Analysis))
	}

}

func (h *Handler) requestAnalysis(r *http.Request, d downloader.Downloader, id, mediaType string, gotF, gotN bool) (a structs.Analysis, err error) {

	filename, reader, err := d.Download(id, mediaType, r)
	if err != nil {
		log.Println("HandleAnalysis: error while downloading file:", err)
		return
	}

	rBytes, err := ioutil.ReadAll(reader)
	_ = reader.Close()
	if err != nil {
		log.Println("HandleAnalysis: error while reading file body:", err)
		return
	}

	//
	fChan := make(chan structs.Analysis)
	fReader := bytes.NewReader(rBytes)
	nChan := make(chan structs.Analysis)
	nReader := bytes.NewReader(rBytes)

	if !gotF {

		go func() {

			f, err := h.performFingerprint(fReader, id, filename, mediaType)
			if err != nil {
				fChan <- structs.Analysis{FingerprintErrorString: err.Error()}
			} else {
				fChan <- structs.Analysis{Fingerprint: f}
			}

		}()

	}

	if !gotN {

		go func() {

			n, err := h.performNSFW(nReader, id, filename, mediaType)
			if err != nil {
				nChan <- structs.Analysis{NSFWErrorString: err.Error()}
			} else {
				nChan <- structs.Analysis{NSFW: n}
			}

		}()

	}

	for !gotF || !gotN {

		select {
		case fAnalysis := <-fChan:
			a.Fingerprint = fAnalysis.Fingerprint
			a.FingerprintErrorString = fAnalysis.FingerprintErrorString
			gotF = true
		case nAnalysis := <-nChan:
			a.NSFW = nAnalysis.NSFW
			a.NSFWErrorString = nAnalysis.NSFWErrorString
			gotN = true
		}

	}

	return a, nil

}

func checkCacheForAnalysis(id string) (fHit, nHit bool, a structs.Analysis) {

	f, err := cache.GetFingerprint(id)
	if err == nil {
		fHit = true
		a.Fingerprint = f
	}

	n, err := cache.GetNSFW(id)
	if err == nil {
		a.NSFW = n
		nHit = true
	}

	return

}

func encodeAnalysisResponse(w http.ResponseWriter, a structs.Analysis) {

	err := json.NewEncoder(w).Encode(a)
	if err != nil {
		log.Println("encodeAnalysisResponse: unable to send response:", err)
	}

}
