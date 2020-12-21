package webserver

import (
	"net/http"
	"strings"
	"time"
)

const (
	MessageTypeDefault   = "message"
	MessageTypeHeartbeat = "heartbeat"
	MessageTypeService   = "service"
)

type Messenger interface {
	Inject(msgType string, msgValue string)
}

func InjectHeartbeat(r *http.Request, m Messenger) {

	for {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(10 * time.Second):
			m.Inject("heartbeat", "ping")
		}
	}
}

func GetVisitorAddress(r *http.Request) string {

	var addr string

	forwarded := r.Header.Get("X-FORWARDED-FOR")

	if forwarded != "" {
		addr = forwarded
	} else {
		addr = r.RemoteAddr
	}

	ip := strings.Split(addr, ":")

	return ip[0]
}
