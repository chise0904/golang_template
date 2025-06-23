package trace

type ctxKey string

const (
	ctxKeyXRequestID          ctxKey = "x-request-id"
	grpcMetadataXRequestIDKey string = "x-request-id"
)

const (
	ctxKeyXTime     ctxKey = "x-time"
	metadataTimeKey        = "x-time"
)
