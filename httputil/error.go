package httputil

import (
	"context"
	"net/http"

	"github.com/okunix/ambicertd/logutil"
)

func HttpError(ctx context.Context, w http.ResponseWriter, err error, statusCode int) {
	if loggerCtx, ok := ctx.Value(logutil.LoggerContextKey).(*logutil.LoggerContext); ok {
		loggerCtx.Error = err
	}
	http.Error(w, err.Error(), statusCode)
}
