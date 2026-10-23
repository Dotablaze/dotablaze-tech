package main

import (
	"github.com/bwmarrin/discordgo"
	"libs/go/meow-bot/feature/handler"
	"libs/go/meow-bot/util"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	logger         = slog.Default()
	botToken       = os.Getenv("DISCORD_BOT_TOKEN")
	allowedChannel = os.Getenv("ALLOWED_CHANNEL_ID")
)

func main() {
	logger.Info("Booting up Meow bot...")

	rand.Seed(time.Now().UnixNano())
	util.InitEmojis(logger)

	if allowedChannel != "" {
		logger.Info("Channel restriction enabled", "channelID", allowedChannel)
	}

	if botToken == "" {
		logger.Error("DISCORD_BOT_TOKEN not set")
		os.Exit(1)
	}

	sess, err := discordgo.New("Bot " + botToken)
	if err != nil {
		logger.Error("Failed to create Discord session", "error", err)
		os.Exit(1)
	}

	sess.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	sess.AddHandler(handler.MessageHandler(logger, allowedChannel))
	sess.AddHandler(handler.CommandHandler(logger))

	err = sess.Open()
	if err != nil {
		logger.Error("Failed to open Discord session", "error", err)
		os.Exit(1)
	}
	defer func(sess *discordgo.Session) {
		err := sess.Close()
		if err != nil {
			logger.Error("Failed to close Discord session", "error", err)
		}
	}(sess)

	_, err = sess.ApplicationCommandCreate(sess.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "meowcount",
		Description: "Check the current meow count",
	})
	if err != nil {
		logger.Error("Failed to register /meowcount", "error", err)
	}

	_, err = sess.ApplicationCommandCreate(sess.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "highscore",
		Description: "Check the highest meow streak",
	})
	if err != nil {
		logger.Error("Failed to register /highscore", "error", err)
	}

	logger.Info("🐱 Meow bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	logger.Info("👋 Shutting down Meow bot.")
}
