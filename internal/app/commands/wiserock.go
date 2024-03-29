package commands

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func wiserock(meta commandeer.Meta, command string, arguments []string) bool {
	var good_advice = []string{
		"If opportunity does not knock, you can always build a door!",
		"Pain doesn't last forever.",
		"Life is short, so smile while you still have teeth.",
		"Follow your heart, but remember to take your brain with you",
		"Do something today that you'll thank yourself for tomorrow!",
	}
	var ok_advice = []string{
		"Not all advice is good advice.",
		"Never make eye contact while eating a banana.",
		"It is not wise to listen to advice from a rock.",
		"Your problems will catch up to you eventually.",
		"There are approximately 238,900 miles between the Earth and the Moon.",
	}
	var bad_advice = []string{
		"When life gives you lemons... Start a LEMONADE stand in front of a TRAIN STATION and give away complimentary TRAIN PASSES with each purchase of LEMONADE.",
		"Don't be ugly.",
		"Yell ALL THE TIME!!!!",
		"Listen to everyone's advice.",
		"Anything is okay as long as you're not caught!",
	}
	rand.Seed(time.Now().Unix())
	var adviceType = []string{"good", "ok", "bad"}[rand.Intn(3)]
	var advices = map[string][]string{
		"good": good_advice,
		"ok":   ok_advice,
		"bad":  bad_advice,
	}

	if len(arguments) > 0 {
		if _, ok := advices[strings.ToLower(arguments[0])]; ok {
			adviceType = strings.ToLower(arguments[0])
		}
	}

	image := meta.DialogueGenerator.GenerateDialogue(advices[adviceType][rand.Intn(len(advices[adviceType]))], meta.DialogueGenerator.NormalFont, 500, 140)
	var buf bytes.Buffer
	err := png.Encode(&buf, image)
	if err != nil {
		logging.Warn("Failed to encode buffer into png", err)
		return false
	}

	rockicon, err := ioutil.ReadFile(path.Join(meta.Settings[settings.ASSETS], "/images/wiserock.png"))
	if err != nil {
		logging.Warn("Failed to read file of wiserock.png", err)
		return false
	}

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL:    "attachment://wisdom.png",
			Width:  500,
			Height: 140,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Wise Rock",
			IconURL: "attachment://icon.png",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v advice", adviceType),
		},
	}

	_, err = meta.Session.ChannelMessageSendComplex(meta.Message.ChannelID, &discordgo.MessageSend{
		Embed: embed,
		Files: []*discordgo.File{
			{
				Name:        "wisdom.png",
				ContentType: "image/png",
				Reader:      bytes.NewReader(buf.Bytes()),
			},
			{
				Name:        "icon.png",
				ContentType: "image/png",
				Reader:      bytes.NewReader(rockicon),
			},
		},
	})
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
