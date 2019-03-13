package main

import (
	"fmt"
	audio "github.com/AlmightyFloppyFish/highlevel-discordgo-opus"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"time"
)

var manage = make(chan audio.AudioAction)

func main() {
	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		fmt.Println("Please provide a discord token by setting the TOKEN environment variable")
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	if dg.Open() != nil {
		panic("Could not open discord session")
	}
	defer dg.Close()

	dg.AddHandler(onMessage)

	select {}
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	if len(words) != 2 || []byte(m.Content)[0] != byte('!') {
		return
	}

	var (
		err error
		vc  *discordgo.VoiceConnection
	)
	switch words[0] {
	case "!link":
		vc, err = audio.JoinUserVoiceChannel(m.Author.ID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel")
			return
		}
		err = audio.AudioFromYoutubeLink(words[1], vc, manage)
	case "!search":
		vc, err = audio.JoinUserVoiceChannel(m.Author.ID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel")
			return
		}
		err = audio.AudioFromYoutubeSearch(words[1], vc, manage)
	case "!oof":
		vc, err = audio.JoinUserVoiceChannel(m.Author.ID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel")
			return
		}
		err = audio.AudioFromFile("assets/oof.mp3", vc, manage)
	}
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to play audio: "+err.Error())
		return
	}

	// Wait for 10 seconds then stops audio
	time.Sleep(10 * time.Second)
	vc.Disconnect()

	/*
		// Alternatively stop without disconnecting from voice channel
		time.Sleep(10 * time.Second)
		manage <- audio.AudioStop

	*/
}
