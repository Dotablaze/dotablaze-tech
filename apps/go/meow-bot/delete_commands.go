package main

import (
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	logger := slog.Default()
	token := os.Getenv("DISCORD_BOT_TOKEN")

	sess, _ := discordgo.New("Bot " + token)
	sess.Open()
	defer sess.Close()

	appID := os.Getenv("DISCORD_APP_ID") // Your application/client ID

	cmds, err := sess.ApplicationCommands(appID, "")
	if err != nil {
		logger.Error("Failed to fetch commands", "error", err)
		os.Exit(1)
	}

	for _, cmd := range cmds {
		err := sess.ApplicationCommandDelete(appID, "", cmd.ID)
		if err != nil {
			logger.Error("Failed to delete", "name", cmd.Name, "error", err)
		} else {
			logger.Info("Deleted", "name", cmd.Name)
		}
	}
}
