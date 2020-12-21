package webserver

import (
	"net/http"
	"time"
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
