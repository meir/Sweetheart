package commandeer

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/meir/Sweetheart/internal/pkg/logging"
	"github.com/meir/Sweetheart/internal/pkg/meta"
	"github.com/meir/Sweetheart/internal/pkg/settings"
)

type Commandeer struct {
	prefix   string
	commands map[string]struct {
		cmd   Command
		arg   Arguments
		usage string
	}
	FailedArguments Command
	meta            *meta.Meta
	Maintenance     bool
}

type Meta struct {
	*meta.Meta
	User    *discordgo.User
	Message *discordgo.Message
	Usage   string
}

type Command func(meta Meta, command string, arguments []string) bool

func NewCommandeer(prefix string, meta *meta.Meta) *Commandeer {
	return &Commandeer{
		prefix: prefix,
		commands: map[string]struct {
			cmd   Command
			arg   Arguments
			usage string
		}{},
		FailedArguments: nil,
		meta:            meta,
	}
}

func (c *Commandeer) Apply(command string, cmd Command, arguments Arguments, usage string, description string) {
	c.commands[strings.ToLower(command)] = struct {
		cmd   Command
		arg   Arguments
		usage string
	}{
		cmd:   cmd,
		arg:   arguments,
		usage: usage,
	}
	if !arguments.Hidden {
		c.meta.Commands[fmt.Sprintf("%v%v %v", c.meta.Settings[settings.PREFIX], strings.ToLower(command), usage)] = description
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
		var done = false

		if cmd, ok := c.commands[strings.ToLower(command)]; ok {
			if cmd.arg.Any {
				goto accepted
			}
			if !(cmd.arg.Min <= 0 || cmd.arg.Min <= len(args)) {
				goto failed
			}
			if !(cmd.arg.Max <= 0 || cmd.arg.Max >= len(args)) {
				goto failed
			}
			for _, amount := range cmd.arg.Amounts {
				if len(args) == amount {
					goto accepted
				}
			}
			if len(cmd.arg.Amounts) > 0 {
				goto failed
			}
			if c.meta.Settings[settings.MAINTENANCE] == "true" && !c.meta.IsAdmin(msg.Author.ID) {
				_, err := session.ChannelMessageSend(msg.ChannelID, "Sorry, but im currently in maintenance mode.")
				if err != nil {
					logging.Warn("Failed to send message", err)
					return
				}
				return
			}

		accepted:
			done = cmd.cmd(Meta{
				c.meta, msg.Author, msg, cmd.usage,
			}, command, args)
			if !done {
				goto failed
			}
			return

		failed:

			// Failed argument checks
			if c.FailedArguments != nil {
				c.FailedArguments(Meta{
					c.meta, msg.Author, msg, cmd.usage,
				}, command, args)
			} else {
				_, err := session.ChannelMessageSend(msg.ChannelID, "[E] failed argument check; no argument failcheck function setup!")
				if err != nil {
					logging.Warn("Failed to send message", err)
					return
				}
			}
		}
	}
}
