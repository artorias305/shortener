package handler

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	RespondJson(w, r, http.StatusOK, "status: ok")
}
