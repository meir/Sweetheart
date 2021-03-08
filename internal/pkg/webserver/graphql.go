package webserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/graphql-go/graphql"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ImageID       string   `bson:"image_id"`
	About         string   `bson:"about"`
	Description   string   `bson:"description"`
	FavoriteColor int      `bson:"favorite_color"`
	Socials       []Social `bson:"socials"`
	Timezone      string   `bson:"timezone"`
	Country       string   `bson:"country"`

	Gender    string `bson:"gender"`
	Pronouns  string `bson:"pronouns"`
	Sexuality string `bson:"sexuality"`
}

type DiscordDetails struct {
	ID            string `json:"id" bson:"id"`
	Username      string `json:"username" bson:"username"`
	Avatar        string `json:"avatar" bson:"avatar"`
	Discriminator string `json:"discriminator" bson:"discriminator"`
	Profile       User   `json:"user,omitempty" bson:"profile"`
}

type Social struct {
	Name   string `bson:"name"`
	Handle string `bson:"handle"`
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

func (ws *Webserver) identity() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Identity",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"avatar": &graphql.Field{
				Type: graphql.String,
			},
			"discriminator": &graphql.Field{
				Type: graphql.String,
			},
			"picture": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					details := p.Source.(DiscordDetails)
					return fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%v", details.ID, details.Avatar), nil
				},
			},
		},
	})
}

func (ws *Webserver) profile() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Profile",
		Fields: graphql.Fields{
			"image_id": &graphql.Field{
				Type: graphql.String,
			},
			"about": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"favorite_color": &graphql.Field{
				Type: graphql.String,
			},
			// "socials": &graphql.Field{
			// 	Type: graphql.String,
			// },
			"timezone": &graphql.Field{
				Type: graphql.String,
			},
			"country": &graphql.Field{
				Type: graphql.String,
			},
			"gender": &graphql.Field{
				Type: graphql.String,
			},
			"pronouns": &graphql.Field{
				Type: graphql.String,
			},
			"sexuality": &graphql.Field{
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
		"identity": &graphql.Field{
			Type: ws.identity(),
			Args: graphql.FieldConfigArgument{
				"session": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				session := p.Args["session"].(string)
				jsonString, err := base64.StdEncoding.DecodeString(session)
				if err != nil {
					return nil, err
				}

				var details = DiscordDetails{}
				err = json.Unmarshal([]byte(jsonString), &details)
				if err != nil {
					return nil, err
				}

				database := ws.Meta.Database.Database("users")
				collection := database.Collection("users")
				upsert := true
				details.Profile = User{
					About:         "I'm absolutely amazing!",
					Description:   fmt.Sprintf("Hi, i'm %v and i'm absolutely amazing abviously!", details.Username),
					FavoriteColor: 0xffffff,
					Socials:       []Social{},
					Timezone:      "CET",
					Country:       "Netherlands",

					Gender:    "???",
					Pronouns:  "???/???/???",
					Sexuality: "???",
				}
				_, err = collection.UpdateOne(context.Background(), bson.M{
					"ID": details.ID,
				}, bson.M{
					"$set": bson.M{
						"username":      details.Username,
						"avatar":        details.Avatar,
						"discriminator": details.Discriminator,
					},
					"$setOnInsert": bson.M{
						"id":      details.ID,
						"profile": details.Profile,
					},
				}, &options.UpdateOptions{
					Upsert: &upsert,
				})
				return details, err
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
