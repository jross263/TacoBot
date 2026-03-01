package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Handler interface {
	ID() string
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type ApplicationCommand interface {
	Handler
	Definition() *discordgo.ApplicationCommand
}

var handlers = make(map[string]Handler)

func Add(h Handler) {
	_, ok := handlers[h.ID()]
	if ok {
		panic(fmt.Sprintf("handler with ID %q is already registered", h.ID()))
	}
	handlers[h.ID()] = h
}

func Get(id string) (Handler, bool) {
	h, ok := handlers[id]
	return h, ok
}

func AllCommands() []ApplicationCommand {
	var cmds []ApplicationCommand
	for _, h := range handlers {
		if cmd, ok := h.(ApplicationCommand); ok {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}
