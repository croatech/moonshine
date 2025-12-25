package api

import (
	"bytes"
	"io"
	"log"

	"github.com/labstack/echo/v4"
)

func cacheRequestBody() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Body != nil {
				bodyBytes, err := io.ReadAll(c.Request().Body)
				if err == nil {
					c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					c.Set("body", bodyBytes)
					// Log request body for debugging
					if len(bodyBytes) > 0 {
						log.Printf("[GraphQL Request Body] %s", string(bodyBytes))
					}
				}
			}
			return next(c)
		}
	}
}

func getCachedBody(c echo.Context) []byte {
	if body, ok := c.Get("body").([]byte); ok {
		return body
	}
	return nil
}

