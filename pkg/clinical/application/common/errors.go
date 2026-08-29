package common

import "errors"

// ErrUnknownTopic is returned when an event arrives on a topic this service does
// not consume. The push handler treats it as a bad request; the NATS subscriber
// acknowledges and ignores it.
var ErrUnknownTopic = errors.New("unknown topic")
