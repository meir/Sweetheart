package webserver

import (
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
	}

	root := graphql.NewObject(graphql.ObjectConfig{
		Name:   "query",
		Fields: queryFields,
	})
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: root}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logging.Warn("failed to create new schema, error: %v", err)
	}
	return &schema
}
