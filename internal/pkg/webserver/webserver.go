package webserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/meta"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type Webserver struct {
	Meta   *meta.Meta
	Schema *graphql.Schema
}

func NewWebserver(meta *meta.Meta) *Webserver {
	return &Webserver{meta, nil}
}

func (ws *Webserver) Start() {
	defer func() {
		ws.Meta.Status["Webserver"] = false
	}()

	ws.Schema = ws.schema()

	ws.Meta.Status["Webserver"] = true

	err := http.ListenAndServe(fmt.Sprintf(":%v", ws.Meta.Settings[settings.PORT]), http.HandlerFunc(ws.handler))
	if err != nil {
		panic(err)
	}
}

func (ws *Webserver) handler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api") {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
			logging.Warn("failed to read body of an API request.", err)
			return
		}
		params := graphql.Params{
			Schema:        *ws.Schema,
			RequestString: string(body),
		}
		r := graphql.Do(params)
		if len(r.Errors) > 0 {
			logging.Warn("failed to execute graphql operation, errors:", r.Errors)
		}
		response, err := json.Marshal(r)
		if err != nil {
			logging.Warn("failed to parse graphql response to JSON.", err)
			return
		}
		w.WriteHeader(200)
		w.Write(response)
	} else if r.URL.Path == "/heartbeat" {
		w.WriteHeader(200)
		w.Write([]byte("Meow? (Waiting for something to happen?)"))
	} else {
		http.ServeFile(w, r, path.Join(ws.Meta.Settings[settings.ASSETS], "web", r.URL.Path))
	}
}
