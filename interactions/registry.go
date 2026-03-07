package interactions

import (
	"slices"
	"tacobot/util"

	"github.com/bwmarrin/discordgo"
)

type HandleFunc func(*discordgo.Session, *discordgo.InteractionCreate, map[string]string) error

type Handler struct {
	ID     string
	Handle HandleFunc
}

type Command struct {
	Definition   *discordgo.ApplicationCommand
	Handle       HandleFunc
	Interactions []Handler
}

type registry struct {
	handlers *util.Cache[string, HandleFunc]
	commands []*discordgo.ApplicationCommand
}

var reg = &registry{handlers: util.NewCache[string, HandleFunc]("handlers")}

func Add(cmd Command) {
	id := cmd.Definition.Name

	reg.handlers.AddOnce(id, cmd.Handle)
	reg.commands = append(reg.commands, cmd.Definition)
	for _, i := range cmd.Interactions {
		reg.handlers.AddOnce(i.ID, i.Handle)
	}
}

func Get(id string) (HandleFunc, bool) {
	return reg.handlers.Get(id)
}

func AllCommands() []*discordgo.ApplicationCommand {
	return slices.Clone(reg.commands)
}
