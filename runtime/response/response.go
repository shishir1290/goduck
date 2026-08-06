package response

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	writer *Recorder
}

func New(
	w *Recorder,
) *Response {

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

func (r *Response) StatusCode() int {

	return r.writer.Status
}

func (r *Response) Size() int {

	return r.writer.Size
}