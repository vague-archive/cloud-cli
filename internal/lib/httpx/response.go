package httpx

import (
	"encoding/json"
	"net/http"
)

func RespondOk(value any, w http.ResponseWriter) {
	Respond(http.StatusOK, value, w)
}

func RespondAccepted(value any, w http.ResponseWriter) {
	Respond(http.StatusAccepted, value, w)
}

func RespondBadRequest(value any, w http.ResponseWriter) {
	Respond(http.StatusBadRequest, value, w)
}

func Respond(status int, value any, w http.ResponseWriter) {
	switch v := value.(type) {
	case string:
		w.Header().Add(HeaderContentType, ContentTypeText)
		w.WriteHeader(status)
		w.Write([]byte(v))
	default:
		w.Header().Add(HeaderContentType, ContentTypeJSON)
		bytes, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(status)
		w.Write(bytes)
	}
}
