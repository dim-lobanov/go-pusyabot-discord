package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Token          string
	nextImageIndex int = 0
	randomPerm     []int
)

const imagesDir = "images"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	botSession, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("Error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	botSession.AddHandler(messageCreate)

	// We only care about receiving message events.
	botSession.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = botSession.Open()
	if err != nil {
		log.Println("Error opening connection,", err)
		log.Println("Token: ", Token)
		return
	}
	defer botSession.Close()

	fmt.Println("Bot is now running!")
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if message.Author.ID == session.State.User.ID {
		return
	}

	if message.Content == "!pusya" || message.Content == "!пуся" || strings.ToLower(message.Content) == "pusyabot pusya me" {
		go pusyaCommand(session, message)
	}
}

func pusyaCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	image := getImageFile(imagesDir)
	defer image.Close()

	_, err := session.ChannelFileSend(message.ChannelID, "pusya.png", image)
	if err != nil {
		log.Println(err)
	}
}

func getImageFile(imagesDir string) *os.File {
	fileInfoList, err := os.ReadDir(imagesDir)
	if err != nil {
		log.Fatal(err)
	}

	n := getNextRandomInt(len(fileInfoList))
	reader, err := os.Open(path.Join(imagesDir, fileInfoList[n].Name()))
	if err != nil {
		log.Println(err)
	}

	return reader
}

func getNextRandomInt(upperBound int) int {
	if len(randomPerm) == 0 {
		randomPerm = rand.Perm(upperBound)
		log.Println("Generating next permutation: ", randomPerm)
	}

	if nextImageIndex >= len(randomPerm) {
		randomPerm = rand.Perm(upperBound)
		nextImageIndex = 0
	}

	randomInt := randomPerm[nextImageIndex]
	nextImageIndex += 1
	return randomInt
}
