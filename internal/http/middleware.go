package http

import (
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//AddRequestIDToContext is a middleware that create a logger with a request id
func AddRequestIDToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			c.Set("request_id", requestID)
			return next(c)
		}
	}
}

//AddLoggerToContext is a middleware that create a logger with a request id
func AddLoggerToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			logger := log.WithFields(log.Fields{
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
			})
			c.Set("logger", logger)
			return next(c)
		}
	}
}

//Tracing adds tracing capabilities to REST endpoints
func Tracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//no tracing for static files
			if c.Path() == "/*" || c.Path() == "/metrics" {
				return next(c)
			}

			tracer := opentracing.GlobalTracer()
			req := c.Request()
			opName := "HTTP " + req.Method + " URL: " + c.Path()

			var span opentracing.Span
			if ctx, err := tracer.Extract(opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(req.Header)); err != nil {
				span = tracer.StartSpan(opName)
			} else {
				span = tracer.StartSpan(opName, ext.RPCServerOption(ctx))
			}

			ext.HTTPMethod.Set(span, req.Method)
			ext.HTTPUrl.Set(span, req.URL.String())
			ext.Component.Set(span, "rest")

			req = req.WithContext(opentracing.ContextWithSpan(req.Context(), span))
			c.SetRequest(req)

			c.Set("span", span)

			defer func() {
				status := c.Response().Status
				committed := c.Response().Committed
				ext.HTTPStatusCode.Set(span, uint16(status))
				if status >= http.StatusInternalServerError || !committed {
					ext.Error.Set(span, true)
				}
				span.Finish()
			}()

			return next(c)
		}
	}
}
