package ncore

const namespace = "nbs-go/nucleo/ncore"
const wrappedErrorFmt = "\n  > %w"

type Core struct {
	Manifest    Manifest
	Environment Environment
	WorkDir     string
	NodeId      string
	Responses   *ResponseMap
}
