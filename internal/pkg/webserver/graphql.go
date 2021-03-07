package webserver

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
