package transport

import (
	"fmt"
	"strings"
	"time"

	"github.com/ghettovoice/gosip/core"
	"github.com/ghettovoice/gosip/log"
)

const (
	netErrRetryTime = 5 * time.Second
	socketTtl       = time.Hour
)

// Protocol implements network specific transport features.
type Protocol interface {
	log.WithLogger
	Network() string
	Reliable() bool
	Streamed() bool
	Listen(target *Target) error
	Send(target *Target, msg core.Message) error
	String() string
}

type protocol struct {
	log      log.Logger
	network  string
	reliable bool
	streamed bool
}

func (pr *protocol) SetLog(logger log.Logger) {
	pr.log = logger.WithFields(map[string]interface{}{
		"protocol": pr.String(),
	})
}

func (pr *protocol) Log() log.Logger {
	return pr.log
}

func (pr *protocol) String() string {
	var name, network string
	if pr == nil {
		name = "<nil>"
		network = ""
	} else {
		name = fmt.Sprintf("%p", pr)
		network = pr.Network() + " "
	}

	return fmt.Sprintf("%sprotocol %p", network, name)
}

func (pr *protocol) Network() string {
	return strings.ToUpper(pr.network)
}

func (pr *protocol) Reliable() bool {
	return pr.reliable
}

func (pr *protocol) Streamed() bool {
	return pr.streamed
}