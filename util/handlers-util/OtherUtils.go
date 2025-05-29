package handlersUtil

import "net/http"

func writeResponse(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
