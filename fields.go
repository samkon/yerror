package yerror

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

var (
	CodeBadRequest           = zap.Int("code", http.StatusBadRequest)
	CodeUnauthorized         = zap.Int("code", http.StatusUnauthorized)
	CodeForbidden            = zap.Int("code", http.StatusForbidden)
	CodeNotFound             = zap.Int("code", http.StatusNotFound)
	CodeMethodNotAllowed     = zap.Int("code", http.StatusMethodNotAllowed)
	CodeNotAcceptable        = zap.Int("code", http.StatusNotAcceptable)
	CodeRequestTimeout       = zap.Int("code", http.StatusRequestTimeout)
	CodePreconditionFailed   = zap.Int("code", http.StatusPreconditionFailed)
	CodeUnsupportedMediaType = zap.Int("code", http.StatusUnsupportedMediaType)
	CodeTooManyRequests      = zap.Int("code", http.StatusTooManyRequests)
	CodeInternalServerError  = zap.Int("code", http.StatusInternalServerError)
	CodeNotImplemented       = zap.Int("code", http.StatusNotImplemented)
	CodeBadGateway           = zap.Int("code", http.StatusBadGateway)
	CodeServiceUnavailable   = zap.Int("code", http.StatusServiceUnavailable)
	CodeGatewayTimeout       = zap.Int("code", http.StatusGatewayTimeout)
)

func Code(code int) zap.Field {
	return zap.Int("code", code)
}

// JSON Returns raw encoded JSON value if given bytes is serializable json. Otherwise, returns byte string.
func JSON(key string, value []byte) zap.Field {
	if json.Valid(value) {
		return zap.Any(key, json.RawMessage(value))
	}
	return zap.ByteString(key, value)
}
