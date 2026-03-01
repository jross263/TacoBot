package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"tacobot/commands"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type config struct {
	appId   string
	guildId string
	token   string
}

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	config, err := loadConfig()
	if err != nil {
		return err
	}

	session, err := discordgo.New("Bot " + config.token)
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	slog.Info("Attempting to start bot session...")

	err = session.Open()
	if err != nil {
		return fmt.Errorf("opening session: %w", err)
	}

	slog.Info("TacoBot running!")

	registerCommandHandlers(session)
	registeredCommands := registerCommands(session, config)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	removeCommands(session, registeredCommands, config)

	slog.Info("Attempting to stop bot session...")

	err = session.Close()
	if err != nil {
		return fmt.Errorf("closing session: %w", err)
	}

	slog.Info("TacoBot stopped!")

	return nil
}

func loadConfig() (config, error) {
	appId := os.Getenv("DISCORD_APP_ID")
	guildId := os.Getenv("DISCORD_GUILD_ID")
	token := os.Getenv("DISCORD_TOKEN")

	var missing []string
	if strings.TrimSpace(appId) == "" {
		missing = append(missing, "DISCORD_APP_ID")
	}
	if strings.TrimSpace(guildId) == "" {
		missing = append(missing, "DISCORD_GUILD_ID")
	}
	if strings.TrimSpace(token) == "" {
		missing = append(missing, "DISCORD_TOKEN")
	}

	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing requiredx environment variables: %s", strings.Join(missing, ", "))
	}
	return config{appId: appId, guildId: guildId, token: token}, nil
}

func registerCommandHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("Recieved command", "Name", i.ApplicationCommandData().Name)
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if cmd, ok := commands.Get(i.ApplicationCommandData().Name); ok {
				slog.Info("Executing command", "Name", cmd.Definition().Name)
				err := cmd.Handler(s, i)
				if err != nil {
					slog.Error("Error executing command handler", "Name", i.ApplicationCommandData().Name)
				}
			}
		case discordgo.InteractionMessageComponent:
			if cmd, ok := commands.Get(i.MessageComponentData().CustomID); ok {
				if cc, ok := cmd.(commands.ComponentCommand); ok {
					err := cc.ComponentHandler(s, i)
					if err != nil {
						slog.Error("Error executing component handler", "CustomID", i.MessageComponentData().CustomID)
					}
				}
			}
		}
	})
}

func registerCommands(s *discordgo.Session, config config) []*discordgo.ApplicationCommand {
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(commands.All()))

	for _, cmd := range commands.All() {
		slog.Info("Attempting to add command", "Name", cmd.Definition().Name)
		createdCmd, err := s.ApplicationCommandCreate(config.appId, config.guildId, cmd.Definition())
		if err != nil {
			slog.Error("Error adding command", "Name", cmd.Definition().Name, "err", err)
			continue
		}
		slog.Info("Added command", "Name", cmd.Definition().Name)
		registeredCommands = append(registeredCommands, createdCmd)
	}

	return registeredCommands
}

func removeCommands(s *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand, config config) {
	cmds, err := s.ApplicationCommands(config.appId, config.guildId)

	if err != nil {
		slog.Warn("Error retrieving commands, defaulting to registered list")
		cmds = registeredCommands
	}

	for _, cmd := range cmds {
		slog.Info("Attempting to remove command", "Name", cmd.Name, "ID", cmd.ID)
		err := s.ApplicationCommandDelete(config.appId, config.guildId, cmd.ID)
		if err != nil {
			slog.Error("Error trying to delete command", "Name", cmd.Name, "ID", cmd.ID)
			continue
		}
		slog.Info("Removed command", "Name", cmd.Name)
	}
}
