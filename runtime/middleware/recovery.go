package middleware

import (
	"fmt"
	"net/http"

	"github.com/shishir1290/goduck/runtime/context"
	"github.com/shishir1290/goduck/runtime/handler"
)

func Recovery() Middleware {

	return func(next handler.HandlerFunc) handler.HandlerFunc {

		return func(ctx *context.Context) {

			defer func() {

				if err := recover(); err != nil {

					fmt.Println("Recovered Panic:", err)

					ctx.Status(http.StatusInternalServerError)

					ctx.String("Internal Server Error")
				}

			}()

			next(ctx)
		}
	}
}