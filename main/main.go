package main

import (
	"OneNightUltimateWerewolf/main/GD"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token       string
	isGuildIn   map[string]bool
	CardDeckMap map[string]GD.CardDeck
)

func init() {
	isGuildIn = make(map[string]bool)
	CardDeckMap = make(map[string]GD.CardDeck)

	flag.StringVar(&Token, "t", "NzU3OTQzMjQ2NzA0ODY5Mzc4.X2nvpw.f1kQjOdXVjO0ifFVKX6azHwIBQE", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if isGuildIn[m.GuildID] {
		isGuildIn[m.GuildID] = true
		CardDeckMap[m.GuildID] = *GD.NewCardDeck()
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		//cd := CardDeckMap[m.GuildID]
		s.ChannelMessageSend(m.ChannelID, "cd.PopCard()")
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
