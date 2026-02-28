package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")

	if strings.TrimSpace(token) == "" {
		log.Fatal("Discord token not found...")
	}

	session, _ := discordgo.New("Bot " + token)

	log.Print("Attempting to start bot session...")

	err := session.Open()

	if err != nil {
		log.Fatal("An error occurred while opening the session: ", err)
	}

	log.Print("TacoBot running!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Print("Attempting to stop bot session...")

	err = session.Close()

	if err != nil {
		log.Fatal("An error occurred while closing the session: ", err)
	}
}
