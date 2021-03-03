package commandeer

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type Commandeer struct {
	prefix   string
	settings map[settings.BotSetting]string
	commands map[string]struct {
		cmd Command
		arg Arguments
	}
	FailedArguments *Command
}

type Meta struct {
	Session  *discordgo.Session
	User     *discordgo.User
	Message  *discordgo.Message
	Settings map[settings.BotSetting]string
}

type Command func(meta Meta, command string, arguments []string) bool

func NewCommandeer(prefix string, st map[settings.BotSetting]string) *Commandeer {
	return &Commandeer{
		prefix:   prefix,
		settings: st,
		commands: map[string]struct {
			cmd Command
			arg Arguments
		}{},
		FailedArguments: nil,
	}
}

func (c *Commandeer) Apply(command string, cmd Command, arguments Arguments) {
	c.commands[strings.ToLower(command)] = struct {
		cmd Command
		arg Arguments
	}{
		cmd: cmd,
		arg: arguments,
	}
}

func (c *Commandeer) Start(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if s.State.User.ID != m.Author.ID {
			c.Run(s, m.Message)
		}
	})
}

func (c *Commandeer) Run(session *discordgo.Session, msg *discordgo.Message) {
	if strings.HasPrefix(msg.Content, c.prefix) {
		message := msg.Content
		arguments := strings.Split(message, " ")
		command := strings.Replace(arguments[0], c.prefix, "", 1)
		args := arguments[1:]

		if cmd, ok := c.commands[strings.ToLower(command)]; ok {
			if cmd.arg.Any {
				goto accepted
			}
			if cmd.arg.Min < 0 || cmd.arg.Min <= len(args) {
				goto accepted
			}
			if cmd.arg.Max < 0 || cmd.arg.Max >= len(args) {
				goto accepted
			}
			if len(cmd.arg.Amounts) == 0 {
				goto accepted
			}
			for _, amount := range cmd.arg.Amounts {
				if len(args) == amount {
					goto accepted
				}
			}

			// Failed argument checks
			if c.FailedArguments != nil {
				(*c.FailedArguments)(Meta{
					session, msg.Author, msg, c.settings,
				}, command, args)
			} else {
				session.ChannelMessageSend(msg.ChannelID, "[E] failed argument check; no argument failcheck function setup!")
			}
			return

		accepted:
			cmd.cmd(Meta{
				session, msg.Author, msg, c.settings,
			}, command, args)
		}
	}
}
