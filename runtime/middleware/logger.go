package middleware

import (
	"fmt"
	"time"

	"github.com/shishir1290/goduck/runtime/context"
	"github.com/shishir1290/goduck/runtime/handler"
)

func Logger() Middleware {

	return func(
		next handler.HandlerFunc,
	) handler.HandlerFunc {

		return func(ctx *context.Context) {

			start := time.Now()

			next(ctx)

			fmt.Printf(
				"%s %-20s %3d %6dB %v\n",
				ctx.Request.Method(),
				ctx.Request.Path(),
				ctx.StatusCode(),
				ctx.ResponseSize(),
				time.Since(start),
			)
		}
	}
}