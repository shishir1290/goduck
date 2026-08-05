package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	writer http.ResponseWriter
}

func New(w http.ResponseWriter) *Response {
	return &Response{
		writer: w,
	}
}

func (r *Response) String(text string) {

	fmt.Fprintln(r.writer, text)

}

func (r *Response) JSON(v any) error {

	r.writer.Header().Set(
		"Content-Type",
		"application/json",
	)

	return json.NewEncoder(r.writer).Encode(v)
}

func (r *Response) Status(code int) {

	r.writer.WriteHeader(code)

}