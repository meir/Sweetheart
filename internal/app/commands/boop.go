package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/meir/Sweetheart/internal/pkg/commandeer"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

func boop(meta commandeer.Meta, command string, arguments []string) bool {
	if meta.Message.Author.ID != "213603030384377856" {
		return true
	}
	images := []string{
		"perfectheart.png",
		"sweetheart.png",
		"mutantheart.png",
	}

	i, err := strconv.Atoi(arguments[0])
	if err != nil {
		return false
	}

	if !(i >= 0 && i < len(images)) {
		return false
	}

	img, err := ioutil.ReadFile(path.Join(meta.Settings[settings.ASSETS], "images", images[i]))
	if err != nil {
		logging.Warn("Failed to read file", err)
		return false
	}
	t := http.DetectContentType(img)
	im := base64.StdEncoding.EncodeToString(img)

	_, err = meta.Session.UserUpdate("", "", "", fmt.Sprintf("data:%s;base64,%s", t, im), "")
	if err != nil {
		logging.Warn("Failed to send message", err)
		return false
	}
	return true
}
