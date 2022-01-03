package pubsub

import "github.com/nbs-go/nlogger"

var logger nlogger.Logger

func init() {
	logger = nlogger.Get()
}
