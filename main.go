package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
	gocord "github.com/nextep-community/gocord"
	"github.com/nextep-community/gocord/bot"
	"github.com/nextep-community/gocord/discord"
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

	r.SlashCommand("/ping", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		return e.CreateMessage(discord.MessageCreate{
			Content: "pong",
			Components: []discord.ContainerComponent{
				discord.ActionRowComponent{
					discord.NewPrimaryButton("Test", "/test"),
				},
			},
		})
	})

	client, err := gocord.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentsAll,
			),
		),
		bot.WithEventListeners(r),
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

var Commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "ping",
		Description: "Replies with pong",
	},
}

func registerCommands(client bot.Client) {
	response, err := client.Rest().SetGlobalCommands(snowflake.MustParse(os.Getenv("APPLICATION_ID")), Commands)

	if err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
		panic(err)
	}

	for _, command := range response {
		slog.Info("registered command", slog.Any("command", command))
	}
}
