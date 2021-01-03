package handler

import (
	"github.com/gorilla/mux"
	"github.com/shitpostingio/analysis-api/api/database"
	"github.com/shitpostingio/analysis-commons/downloader"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Handler contains values needed to dispatch the requests correctly.
type Handler struct {
	NSFWEndpoint        string
	FingerprintEndpoint string
	GibberishEndpoint   string
	Testing             bool
	Downloader          downloader.Downloader
}

const (
	fpKeyHeaderName string = "X-shitposting-key"
	//callerAPIKeyHeaderName   string = "X-caller-bot-apikey"
	//downloadURLHeaderName    string = "X-download-file-url"
	gibberishInputHeaderName string = "Gibberish-Input"
)

// New creates a new handler, given the endpoints where to perform requests.
func New(testing bool, nsfwEndpoint, fingerprintEndpoint, gibberishEndpoint string, d downloader.Downloader) *Handler {
	return &Handler{
		NSFWEndpoint:        nsfwEndpoint,
		FingerprintEndpoint: fingerprintEndpoint,
		GibberishEndpoint:   gibberishEndpoint,
		Testing:             testing,
		Downloader:          d,
	}
}

// Dispatch checks the request for authorization and dispatches it to the
// appropriate handler.
func (h *Handler) Dispatch(w http.ResponseWriter, r *http.Request) {

	log.Debugln("Received request")

	if !h.Testing {
		fpAuthKey := r.Header.Get(fpKeyHeaderName)
		if !database.IsAuthorized(fpAuthKey) {
			log.Println("Received unauthorized token:", fpAuthKey)
			http.Error(
				w,
				http.StatusText(http.StatusForbidden),
				http.StatusForbidden,
			)
			return
		}
	}

	log.Debugln("Authenticated!")
	switch mux.Vars(r)["reqType"] {
	case "complete":
		h.HandleAnalysis(w, r, h.Downloader)
	case "fingerprint":
		h.HandleFingerprint(w, r, h.Downloader)
	case "nsfw":
		h.HandleNSFW(w, r, h.Downloader)
	case "gibberish":
		h.HandleGibberish(w, r)
	default:
		http.Error(w, "no type specified", http.StatusBadRequest)
	}

}
