package response

import "net/http"

type Recorder struct {
	http.ResponseWriter

	Status int

	Size int
}

func NewRecorder(
	w http.ResponseWriter,
) *Recorder {

	return &Recorder{
		ResponseWriter: w,
		Status:         http.StatusOK,
	}
}

func (r *Recorder) WriteHeader(
	status int,
) {

	r.Status = status

	r.ResponseWriter.WriteHeader(status)
}

func (r *Recorder) Write(
	data []byte,
) (int, error) {

	n, err := r.ResponseWriter.Write(data)

	r.Size += n

	return n, err
}