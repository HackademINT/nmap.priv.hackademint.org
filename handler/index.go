package handler

import (
	"net/http"
	"encoding/json"
)

type indexHandler struct {
	jsonResult []byte
}

func NewIndexHandler() (http.Handler, error) {
	var indexContent = map[string]interface{}{
		"author":    "zTeeed",
		"follow_me": "https://github.com/zteeed",
		"paths": []string{
			"/ip",
		},
	}

	jsonResult, err := json.Marshal(indexContent)
	if err != nil {
		return nil, err
	}
	return &indexHandler{jsonResult: jsonResult}, nil

}
func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(h.jsonResult)
}
