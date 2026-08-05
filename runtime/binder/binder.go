package binder

import (
	"encoding/json"
	"net/http"
)

func Bind(
	r *http.Request,
	v any,
) error {

	return json.NewDecoder(r.Body).Decode(v)
}