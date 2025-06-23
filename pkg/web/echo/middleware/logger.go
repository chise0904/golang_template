package middleware

import (
	"bytes"
	"io"
	"net/http/httputil"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func getBody(c echo.Context, maxLogBodySize int) ([]byte, error) {
	reqBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))
	if maxLogBodySize > 0 && len(reqBody) > maxLogBodySize {
		return append(reqBody[:maxLogBodySize], []byte("...")...), nil
	}
	return reqBody, nil
}

// NewAccessLogMiddleware ...
func NewAccessLogMiddleware(inputAndRequestDump bool, maxLogBodySize int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()

			if req.RequestURI == "/ping" {
				return nil
			}

			ctx := req.Context()
			event := log.Ctx(ctx).Debug()

			var reqBody []byte
			var reqDump []byte
			var resBody *bytes.Buffer
			start := time.Now()

			if inputAndRequestDump {
				reqBody, err = io.ReadAll(c.Request().Body)
				if err != nil {
					return err
				}

				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))

				if len(reqBody) > 0 {
					if maxLogBodySize > 0 && len(reqBody) > maxLogBodySize {
						event = event.Str("input", string(reqBody[:maxLogBodySize])+"...")
					} else {
						event = event.Str("input", string(reqBody))
					}
				}
				reqDump, _ = httputil.DumpRequest(req, false)
				event = event.Str("req_dump", string(reqDump))

				resBody = new(bytes.Buffer)
				mw := io.MultiWriter(c.Response().Writer, resBody)
				c.Response().Writer = &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			}

			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()
			latency := stop.Sub(start)
			res := c.Response()

			if inputAndRequestDump {
				respDump := resBody.Bytes()
				if maxLogBodySize > 0 && len(respDump) > maxLogBodySize {
					respDump = append(respDump[:maxLogBodySize], []byte("...")...)
				}
				event = event.Str("resp_dump", string(respDump))
			}

			event = event.Str("host", req.Host).
				Str("uri", req.RequestURI).
				Str("method", req.Method).
				Int("status", res.Status).
				Str("remote_ip", c.RealIP()).
				Str("latency_human", latency.String())

			if err != nil {
				event.Err(err)
			}

			event.Msg("Access Log")

			return nil
		}
	}
}
