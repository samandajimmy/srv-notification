package pubsub

import "github.com/nbs-go/nlogger/v2"

var log nlogger.Logger

func init() {
	log = nlogger.Get()
}
