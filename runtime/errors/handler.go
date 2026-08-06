package errors

import (
	"net/http"

	gcontext "github.com/shishir1290/goduck/runtime/context"
)

func Handle(
	ctx *gcontext.Context,
	err error,
) {

	if e, ok := err.(*HTTPError); ok {

		ctx.Status(e.Status)

		ctx.String(e.Message)

		return
	}

	ctx.Status(http.StatusInternalServerError)

	ctx.String(err.Error())
}