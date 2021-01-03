package handler

import (
	"encoding/json"
	"github.com/shitpostingio/gibberish-microservice/client"
	"log"
	"net/http"
)

// HandleGibberish handles a gibberish request.
func (h *Handler) HandleGibberish(w http.ResponseWriter, r *http.Request) {

	text := r.Header.Get(gibberishInputHeaderName)
	if text == "" {
		http.Error(w, "Gibberish-Input not provided", http.StatusBadRequest)
		return
	}

	data, err := client.PerformRequest(text, h.GibberishEndpoint)
	if err != nil {
		log.Println("HandleGibberish: error while performing gibberish analysis for text ", text, ":", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("HandleGibberish: unable to correctly send response: ", err)
	}

}
