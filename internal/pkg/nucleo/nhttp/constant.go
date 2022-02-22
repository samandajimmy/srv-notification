package nhttp

const (
	// Header keys

	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"

	// Map keys

	MetadataKey        = "metadata"
	HttpStatusRespKey  = "http_status"
	RequestMetadataKey = "request_metadata"
	RequestIdKey       = "requestId"
)

type ContextKey uint8

const (
	_ ContextKey = iota + 1
	RequestIdContextKey
	HttpStatusRespContextKey
	RequestMetadataContextKey
)
