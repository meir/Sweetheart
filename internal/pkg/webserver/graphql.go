package webserver

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/graphql-go/graphql"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type User struct {
	ID      string
	Details DiscordDetails

	ImageID       string
	About         string
	Description   string
	FavoriteColor int
	Socials       []Social
	Timezone      string
	Country       string

	Gender    string
	Pronouns  string
	Sexuality string
}

type DiscordDetails struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
}

type Social struct {
	ID     string
	Name   string
	Handle string
}

func (ws *Webserver) settings() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Settings",
		Fields: graphql.Fields{
			"oauth": &graphql.Field{
				Type: graphql.String,
			},
			"invite": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

func (ws *Webserver) schema() *graphql.Schema {
	queryFields := graphql.Fields{
		"settings": &graphql.Field{
			Type: ws.settings(),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				response := map[string]string{}
				response["oauth"] = ws.Meta.Settings[settings.OAUTH_URL]
				response["invite"] = ws.Meta.Settings[settings.INVITE_URL]
				return response, nil
			},
		},
		"auth": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"code": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				resp, err := http.PostForm("https://discord.com/api/oauth2/token", url.Values{
					"client_id":     {ws.Meta.Settings[settings.CLIENT_ID]},
					"client_secret": {ws.Meta.Settings[settings.CLIENT_SECRET]},
					"grant_type":    {ws.Meta.Settings[settings.GRANT_TYPE]},
					"scope":         {ws.Meta.Settings[settings.SCOPE]},
					"code":          {p.Args["code"].(string)},
					"redirect_uri":  {ws.Meta.Settings[settings.REDIRECT]},
				})
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}

				return base64.StdEncoding.EncodeToString(data), err
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: queryFields}

	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logging.Warn("failed to create new schema, error: %v", err)
	}
	return &schema
}
