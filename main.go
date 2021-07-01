package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/bwmarrin/discordgo"
)

// invite link : https://discord.com/oauth2/authorize?client_id=858208414379409449&scope=bot+applications.commands

// bot list:  https://discord.com/developers/applications

type Counter struct {
	Current string
}

var (
	Token   string
	current int
	Answers []string
)

func main() {

	botToken := os.Getenv("BOT_TOKEN")

	Answers, _ = loadFile("./answers.txt")

	fmt.Println("loaded ", Answers)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

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

	// prefix "!8" found
	if strings.HasPrefix(m.Content, "!8") {
		s.ChannelMessageSend(m.ChannelID, getAnswer())
		return
	}

}

func getAnswer() string {
	rand.Seed(time.Now().Unix())
	return Answers[rand.Intn(len(Answers))]
}

func loadFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.Split(scanner.Text(), "#")[0]
		lines = append(lines, strings.Trim(text, " "))
	}
	return lines, scanner.Err()
}
