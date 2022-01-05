package pubsub

import "github.com/nbs-go/nlogger"

var log nlogger.Logger

func init() {
	log = nlogger.Get()
}
