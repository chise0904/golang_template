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

				// 在 HTTP 處理中，請求的 body 只能被讀取一次。
				// 當你調用 io.ReadAll(c.Request().Body) 後，原始的 body stream 就被消耗完了，後續的中間件或處理函數就無法再讀取 body 內容。
				// 因此，需要將 body 重新放回 buffer 中，以便後續的中間件或處理函數可以再次讀取 body 內容。
				c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))

				if len(reqBody) > 0 {
					if maxLogBodySize > 0 && len(reqBody) > maxLogBodySize {
						event = event.Str("input", string(reqBody[:maxLogBodySize])+"...")
					} else {
						event = event.Str("input", string(reqBody))
					}
				}

				// DumpRequest 用來將 HTTP 請求轉換為可讀的字串格式，讓我詳細解釋：
				// DumpRequest 會產生類似這樣的字串：
				// GET /api/users HTTP/1.1
				// Host: example.com
				// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36
				// Accept: application/json
				// Content-Type: application/json
				// Content-Length: 45

				// {"name": "John", "email": "john@example.com"}
				reqDump, _ = httputil.DumpRequest(req, false)
				event = event.Str("req_dump", string(reqDump))

				resBody = new(bytes.Buffer)
				// io.MultiWriter 創建一個多重寫入器
				// 當資料寫入 mw 時，會同時寫入兩個地方：
				// 原始的 c.Response().Writer（發送給客戶端）
				// resBody buffer（用於日誌記錄）
				mw := io.MultiWriter(c.Response().Writer, resBody)
				c.Response().Writer = &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			}

			if err = next(c); err != nil {
				c.Error(err)
			}
			// ==== next() 後方, 為 response 的 middleware ====

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
