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
				"%s %s %v\n",
				ctx.Request.Method(),
				ctx.Request.Path(),
				time.Since(start),
			)
		}
	}
}