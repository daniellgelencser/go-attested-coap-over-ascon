package client

import (
	"fmt"
	"net"
	"time"

	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/codes"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/message/pool"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net/monitor/inactivity"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/net/responsewriter"
	"github.com/daniellgelencser/go-attested-coap-over-ascon/v3/options/config"
)

const DefaultMTU = 1472

var DefaultConfig = func() Config {
	opts := Config{
		Common: config.NewCommon[*Conn](),
		CreateInactivityMonitor: func() InactivityMonitor {
			return inactivity.NewNilMonitor[*Conn]()
		},
		RequestMonitor: func(*Conn, *pool.Message) (bool, error) {
			return false, nil
		},
		Dialer:                         &net.Dialer{Timeout: time.Second * 3},
		Net:                            "udp",
		TransmissionNStart:             1,
		TransmissionAcknowledgeTimeout: time.Second * 2,
		TransmissionMaxRetransmit:      4,
		GetMID:                         message.GetMID,
		MTU:                            DefaultMTU,
	}
	opts.Handler = func(w *responsewriter.ResponseWriter[*Conn], r *pool.Message) {
		switch r.Code() {
		case codes.POST, codes.PUT, codes.GET, codes.DELETE:
			if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
				opts.Errors(fmt.Errorf("udp client: cannot set response: %w", err))
			}
		}
	}
	return opts
}()

type Config struct {
	config.Common[*Conn]
	CreateInactivityMonitor        CreateInactivityMonitorFunc
	RequestMonitor                 RequestMonitorFunc
	Net                            string
	GetMID                         GetMIDFunc
	Handler                        HandlerFunc
	Dialer                         *net.Dialer
	TransmissionNStart             uint32
	TransmissionAcknowledgeTimeout time.Duration
	TransmissionMaxRetransmit      uint32
	CloseSocket                    bool
	MTU                            uint16
}
