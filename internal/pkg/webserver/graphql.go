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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ImageID       string   `bson:"image_id"`
	About         string   `bson:"about"`
	Description   string   `bson:"description"`
	FavoriteColor int      `json:"favorite_color" bson:"favorite_color"`
	Socials       []Social `bson:"socials"`
	Timezone      int      `bson:"timezone"`
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

func (ws *Webserver) social() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "SocialObject",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"handle": &graphql.Field{
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
					details := p.Source.(*DiscordDetails)
					if details == nil {
						return nil, nil
					}
					return fmt.Sprintf("https://cdn.discordapp.com/avatars/%v/%v", details.ID, details.Avatar), nil
				},
			},
			"profile": &graphql.Field{
				Type: ws.profile(),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					details := p.Source.(*DiscordDetails)
					if details == nil {
						return nil, nil
					}
					collection, err := ws.getCollection("users")
					if err != nil {
						return nil, err
					}
					res := collection.FindOne(context.Background(), bson.M{
						"id": details.ID,
					})
					if res == nil {
						return nil, fmt.Errorf("no profiles found with id of %v", details.ID)
					}
					var profile DiscordDetails
					err = res.Decode(&profile)
					if err != nil {
						return nil, err
					}
					return &profile.Profile, nil
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
				Type: graphql.Int,
			},
			"socials": &graphql.Field{
				Type: graphql.NewList(ws.social()),
			},
			"timezone": &graphql.Field{
				Type: graphql.Int,
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
		"commands": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "Command",
				Fields: graphql.Fields{
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"description": &graphql.Field{
						Type: graphql.String,
					},
				},
			})),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				output := []map[string]string{}
				for k, v := range ws.Meta.Commands {
					output = append(output, map[string]string{
						"name":        k,
						"description": v,
					})
				}
				return output, nil
			},
		},
		"countries": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "Country",
				Fields: graphql.Fields{
					"name": &graphql.Field{
						Type: graphql.String,
					},
					"flag": &graphql.Field{
						Type: graphql.String,
					},
				},
			})),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				var output []struct {
					Name string `json:"name"`
					Flag string `json:"flag"`
				}
				for k, v := range Countries {
					output = append(output, struct {
						Name string `json:"name"`
						Flag string `json:"flag"`
					}{k, v})
				}
				return output, nil
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
				details, err := ws.getUserDetails(session)
				if err != nil {
					return nil, err
				}

				collection, err := ws.getCollection("users")
				if err != nil {
					return nil, err
				}
				upsert := true
				details.Profile = User{
					About:         "I'm absolutely amazing!",
					Description:   fmt.Sprintf("Hi, i'm %v and i'm absolutely amazing obviously!", details.Username),
					FavoriteColor: 0xffffff,
					Socials:       []Social{},
					Timezone:      -60,
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

	mutationFields := graphql.Fields{
		"profile": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"session": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				"about": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				"description": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				"favorite_color": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.Int,
					},
				},
				"socials": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: &graphql.List{
							OfType: graphql.NewInputObject(graphql.InputObjectConfig{
								Name: "Social",
								Fields: graphql.InputObjectConfigFieldMap{
									"name": &graphql.InputObjectFieldConfig{
										Type: graphql.String,
									},
									"handle": &graphql.InputObjectFieldConfig{
										Type: graphql.String,
									},
								},
							}),
						},
					},
				},
				"timezone": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.Int,
					},
				},
				"country": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				"gender": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				"pronouns": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
				"sexuality": &graphql.ArgumentConfig{
					Type: &graphql.NonNull{
						OfType: graphql.String,
					},
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				session := p.Args["session"].(string)
				details, err := ws.getUserDetails(session)
				if err != nil {
					return nil, err
				}

				country := p.Args["country"].(string)
				if _, ok := Countries[country]; !ok {
					return false, fmt.Errorf("%v is not a usable country", country)
				}

				var socials []Social
				var socialsRaw = p.Args["socials"].([]interface{})
				for i := 0; i < len(socialsRaw); i++ {
					if v, ok := socialsRaw[i].(map[string]interface{}); ok {
						if _, ok := v["name"]; !ok {
							continue
						}
						if _, ok := v["handle"]; !ok {
							continue
						}
						socials = append(socials, Social{
							Name:   v["name"].(string),
							Handle: v["handle"].(string),
						})
					}
				}

				profile := User{
					About:         p.Args["about"].(string),
					Description:   p.Args["description"].(string),
					FavoriteColor: p.Args["favorite_color"].(int),
					Timezone:      p.Args["timezone"].(int),
					Socials:       socials,
					Country:       country,
					Gender:        p.Args["gender"].(string),
					Pronouns:      p.Args["pronouns"].(string),
					Sexuality:     p.Args["sexuality"].(string),
				}

				collection, err := ws.getCollection("users")
				if err != nil {
					return false, err
				}

				_, err = collection.UpdateOne(context.Background(), bson.M{
					"id": details.ID,
				}, bson.M{
					"$set": bson.M{
						"profile": profile,
					},
				})
				return true, err
			},
		},
	}
	rootMutation := graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		logging.Warn("failed to create new schema, error: ", err)
	}
	return &schema
}

type Session struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (ws *Webserver) getUserDetails(session string) (*DiscordDetails, error) {
	data, err := base64.StdEncoding.DecodeString(session)
	if err != nil {
		return nil, err
	}

	var ses Session
	err = json.Unmarshal(data, &ses)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/users/@me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("%v %v", ses.TokenType, ses.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	det, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var details DiscordDetails
	err = json.Unmarshal(det, &details)
	if err != nil {
		return nil, err
	}
	return &details, err
}

func (ws *Webserver) getCollection(name string) (*mongo.Collection, error) {
	database := ws.Meta.Database.Database("sweetheart")
	col := database.Collection(name)
	if col == nil {
		err := database.CreateCollection(context.Background(), name)
		if err != nil {
			return nil, err
		}
		col = database.Collection(name)
	}
	return col, nil
}
