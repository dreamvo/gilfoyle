package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func Send(r http.Handler, method, path string, body interface{}) (*httptest.ResponseRecorder, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, path, bytes.NewReader(data))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}
