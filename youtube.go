package hldiscordopus

import (
	"fmt"
	yts "github.com/KeluDiao/gotube/api"
	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"
	"io"
	"net/http"
	"time"
)

// YoutubeLink contains data for one youtube video as audio source
type YoutubeLink struct {
	Title       string
	Author      string
	Description string
	Duration    time.Duration
	stream      io.Reader
}

func (link YoutubeLink) Play(voice *discordgo.VoiceConnection, manage chan AudioAction) error {
	return AudioFromRaw(link.stream, voice, manage)
}

// DCA has some issues finding the stream by itself
func ytCompatibleStreamFrom(v *ytdl.VideoInfo) (io.Reader, error) {
	format := v.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := v.GetDownloadURL(format)
	if err != nil {
		return nil, fmt.Errorf("Unable to get direct audio stream link from youtube for '%s': %s ", v.Title, err.Error())
	}

	resp, err := http.Get(downloadURL.String())
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
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

		reader, err := ytCompatibleStreamFrom(r)
		if err != nil && debug {
			fmt.Println(err)
		}

		videoInfos[i] = YoutubeLink{
			Title:       r.Title,
			Author:      r.Author,
			Description: r.Description,
			Duration:    r.Duration,
			stream:      reader,
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
