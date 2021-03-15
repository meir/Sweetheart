package commands

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

const buttpoems_url = "https://buttpoems.tumblr.com/random"

var buttpoems_regex = regexp.MustCompile(`<img class="post_media_photo image" src="([^"]*)" alt="image">`)

func buttpoem(meta commandeer.Meta, command string, arguments []string) bool {
	req, err := http.NewRequest(http.MethodGet, buttpoems_url, nil)
	if err != nil {
		logging.Warn("Failed to create new request to buttpoems", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.Warn("Failed to call buttpoems", err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Warn("Failed to read body of buttpoems", err)
			return false
		}
		url := resp.Request.URL.String()
		match := buttpoems_regex.FindStringSubmatch(string(d))
		if len(match) >= 2 {
			_, err := meta.Session.ChannelMessageSendEmbed(meta.Message.ChannelID, &discordgo.MessageEmbed{
				URL:   url,
				Type:  discordgo.EmbedTypeImage,
				Title: "ButtPoems",
				Image: &discordgo.MessageEmbedImage{
					URL: match[1],
				},
			})
			if err != nil {
				logging.Warn("Failed to send message", err)
				return false
			}
			return true
		}
	}
	logging.Warn("Something went wrong while getting buttpoems, statuscode:", resp.StatusCode)
	return false
}
