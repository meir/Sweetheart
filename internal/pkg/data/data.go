package data

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
	ID            string            `json:"id" bson:"id"`
	Username      string            `json:"username" bson:"username"`
	Avatar        string            `json:"avatar" bson:"avatar"`
	Discriminator string            `json:"discriminator" bson:"discriminator"`
	Profile       *User             `json:"user,omitempty" bson:"profile"`
	Ranks         map[string]uint64 `json:"ranks,omitempty" bson:"ranks"`
}

type Social struct {
	Name   string `bson:"name"`
	Handle string `bson:"handle"`
}
