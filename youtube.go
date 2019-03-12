package hldiscordopus

import (
	"fmt"
	yts "github.com/KeluDiao/gotube/api"
	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"
	"time"
)

// YoutubeLink contains data for one youtube video as audio source
type YoutubeLink struct {
	Title       string
	Author      string
	Description string
	Duration    time.Duration
	link        string
}

func (link YoutubeLink) Play(voice *discordgo.VoiceConnection, manage chan AudioAction) error {
	return AudioFromYoutubeLink(link.link, voice, manage)
}

func ytCompatibleLinkFrom(v *ytdl.VideoInfo) (string, error) {
	format := v.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := v.GetDownloadURL(format)
	if err != nil {
		return "", fmt.Errorf("Unable to get direct audio stream link from youtube for '%s': %s ", v.Title, err.Error())
	}
	return downloadURL.String(), nil
}

// GetVideosFromSearch
func GetVideosFromSearch(search string) []YoutubeLink {
	IDs, err := getVideoIDsFromSearch(search, 1)
	if err != nil {
		fmt.Println(err)
		return []YoutubeLink{}
	}

	videoInfos := make([]YoutubeLink, len(IDs))
	for i := range IDs {
		r, err := ytdl.GetVideoInfoFromID(IDs[i])
		if err != nil {
			fmt.Println(err)
		}

		comp, _ := ytCompatibleLinkFrom(r)

		videoInfos[i] = YoutubeLink{
			Title:       r.Title,
			Author:      r.Author,
			Description: r.Description,
			Duration:    r.Duration,
			link:        comp,
		}
	}
	return videoInfos
}

func getVideoIDsFromSearch(keywords string, page int) ([]string, error) {

	search, err := yts.GetSearchUrl(keywords, page)
	if err != nil {
		return []string{}, fmt.Errorf("Unable to search video: %s", err.Error())
	}
	IDs, err := yts.GetVideoIdsFromSearch(search)
	if err != nil {
		return []string{}, fmt.Errorf("Unable to get id's from %s: %s", keywords, err.Error())
	}

	return IDs, nil
}
