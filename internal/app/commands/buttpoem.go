package commands

import (
	"errors"
	"fmt"
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

	url, img, err := getButtPoem()
	if err != nil {
		for i := 0; i < 3 && err != nil; i++ {
			url, img, err = getButtPoem()
		}
	}

	_, err = meta.Session.ChannelMessageSendEmbed(meta.Message.ChannelID, &discordgo.MessageEmbed{
		URL:   url,
		Type:  discordgo.EmbedTypeImage,
		Title: "ButtPoems",
		Image: &discordgo.MessageEmbedImage{
			URL: img,
		},
	})
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return false
}

func getButtPoem() (url, img string, err error) {
	req, err := http.NewRequest(http.MethodGet, buttpoems_url, nil)
	if err != nil {
		logging.Warn("Failed to create new request to buttpoems", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.Warn("Failed to call buttpoems", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		var d []byte
		d, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Warn("Failed to read body of buttpoems", err)
			return
		}
		url = resp.Request.URL.String()
		match := buttpoems_regex.FindStringSubmatch(string(d))
		if len(match) >= 2 {
			img = match[1]
			return
		} else {
			logging.Warn("No regex match has been found: ", fmt.Sprint(match))
			err = errors.New("no regex matches")
			return
		}
	}
	err = errors.New(fmt.Sprintf("Something went wrong while getting buttpoems, statuscode: %v", resp.StatusCode))
	return
}
