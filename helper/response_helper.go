package helper

import (
	"encoding/json"
	"net/http"
	"notefan-golang/models/responses"
)

func ResponseJSON(w http.ResponseWriter, response responses.Response) (int, error) {
	bytes, err := json.Marshal(response.GetBody())

	if err != nil {
		ErrorLog(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.GetCode())
	return w.Write(bytes)
}
