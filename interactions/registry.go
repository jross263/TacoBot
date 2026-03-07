package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type HandleFunc func(*discordgo.Session, *discordgo.InteractionCreate) error

type Handler struct {
	ID     string
	Handle HandleFunc
}

type Command struct {
	Definition   *discordgo.ApplicationCommand
	Handle       HandleFunc
	Interactions []Handler
}

var handlers = make(map[string]HandleFunc)
var commands []*discordgo.ApplicationCommand

func Add(cmd Command) {
	id := cmd.Definition.Name
	_, ok := handlers[id]
	if ok {
		panic(fmt.Sprintf("handler with ID %q is already registered", id))
	}
	commands = append(commands, cmd.Definition)
	handlers[id] = cmd.Handle
	for _, i := range cmd.Interactions {
		handlers[i.ID] = i.Handle
	}
}

func Get(id string) (HandleFunc, bool) {
	h, ok := handlers[id]
	return h, ok
}

func AllCommands() []*discordgo.ApplicationCommand {
	return commands
}
