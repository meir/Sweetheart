package commands

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/commandeer"
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
		panic(err)
	}
	b64image := base64.StdEncoding.EncodeToString(buf.Bytes())

	img, err := ioutil.ReadFile(path.Join(meta.Settings[settings.ASSETS], "/images/wiserock.png"))
	if err != nil {
		panic(err)
	}
	b64icon := base64.StdEncoding.EncodeToString(img)

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL:    fmt.Sprintf("data:image/png;base64,%v", b64image),
			Width:  500,
			Height: 140,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Wise Rock",
			IconURL: fmt.Sprintf("data:image/png;base64,%v", b64icon),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%v advice", adviceType),
		},
	}
	meta.Session.ChannelMessageSendEmbed(meta.Message.ChannelID, embed)
	return true
}
