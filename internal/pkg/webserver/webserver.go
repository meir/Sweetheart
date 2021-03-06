package webserver

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/meir/Sweetheart/internal/pkg/meta"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type Webserver struct {
	Meta *meta.Meta
}

func NewWebserver(meta *meta.Meta) *Webserver {
	return &Webserver{meta}
}

func (ws *Webserver) Start() {
	defer func() {
		ws.Meta.Status["Webserver"] = false
	}()
	ws.Meta.Status["Webserver"] = true

	err := http.ListenAndServe(fmt.Sprintf(":%v", ws.Meta.Settings[settings.PORT]), http.HandlerFunc(ws.handler))
	if err != nil {
		panic(err)
	}
}

func (ws *Webserver) handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		w.WriteHeader(200)
		w.Write([]byte("api"))
	} else if r.URL.Path == "/heartbeat" {
		w.WriteHeader(200)
		w.Write([]byte("Meow? (Waiting for something to happen?)"))
	} else {
		http.ServeFile(w, r, path.Join(ws.Meta.Settings[settings.ASSETS], "web", r.URL.Path))
	}
}
