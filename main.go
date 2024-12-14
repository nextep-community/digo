package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	"github.com/nextep-community/digo/commands"
	gocord "github.com/nextep-community/gocord"
	"github.com/nextep-community/gocord/bot"
	"github.com/nextep-community/gocord/discord"
	"github.com/nextep-community/gocord/events"
	"github.com/nextep-community/gocord/gateway"
	"github.com/nextep-community/gocord/handler"
	"github.com/nextep-community/gocord/handler/middleware"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	applicationID := os.Getenv("APPLICATION_ID")
	token := os.Getenv("DISCORD_TOKEN")

	slog.Info("Starting...")
	slog.Info("gocord version", slog.String("version", gocord.Version))
	slog.Info("Envs", slog.Any("DISCORD_TOKEN:", token), slog.Any("APPLICATION_ID:", applicationID))

	r := handler.New()

	r.Use(middleware.Logger)

	r.NotFound(func(event *handler.InteractionEvent) error {
		return event.CreateMessage(discord.MessageCreate{Content: "not found"})
	})

	r.SlashCommand("/ping", commands.HandlePing)

	client, err := gocord.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentsAll,
			),
		),
		bot.WithEventListeners(r),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnMessageCreate: func(event *events.MessageCreate) {
				if event.Message.Author.Bot {
					return
				}

				message := event.Message
				content := strings.TrimPrefix(message.Content, "!")
				parts := strings.Fields(content)
				command := parts[0]

				switch command {
				case "ping":
					_, _ = event.Client().Rest().CreateMessage(message.ChannelID, discord.MessageCreate{
						Content: "Pong!",
					})

				case "hello":
					_, _ = event.Client().Rest().CreateMessage(message.ChannelID, discord.MessageCreate{
						Content: fmt.Sprintf("Olá, %s!", message.Author.Username),
					})

				default:
					_, _ = event.Client().Rest().CreateMessage(message.ChannelID, discord.MessageCreate{
						Content: fmt.Sprintf("Comando desconhecido: `%s`", command),
					})
				}

			},
		}),
	)

	if err != nil {
		slog.Error("error while building bot", slog.Any("err", err))
		panic(err)
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
		panic(err)
	}

	registerCommands(client)

	slog.Info("example is now running. Press CTRL-C to exit.")

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func registerCommands(client bot.Client) {
	response, err := client.Rest().SetGlobalCommands(snowflake.MustParse(os.Getenv("APPLICATION_ID")), commands.Commands)

	if err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
		panic(err)
	}

	for _, command := range response {
		slog.Info("registered command", slog.Any("command", command))
	}
}
