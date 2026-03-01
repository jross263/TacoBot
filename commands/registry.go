package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	Definition() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type ComponentCommand interface {
	Command
	ComponentHandler(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type Registry struct {
	commands map[string]Command
}

var instance = &Registry{
	commands: make(map[string]Command),
}

func Add(cmd Command) {
	instance.commands[cmd.Definition().Name] = cmd
}

func All() map[string]Command {
	return instance.commands
}

func Get(name string) (Command, bool) {
	cmd, ok := instance.commands[name]
	return cmd, ok
}
