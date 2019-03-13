package gdplay

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

// Play audio from the youtube link
func (link YoutubeLink) Play(voice *discordgo.VoiceConnection) (*AudioSession, error) {
	return AudioFromRaw(link.stream, voice)
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

// GetVideosFromSearch gives an array of videos, which you can manage and then call .Play() on
// WARNING: This method of searching can be quite slow. Reason being that I'm avoiding using the youtube API so
// you don't have to add an API key, if you want better speed then get a developer key and use something like
// godoc.org/google.golang.org/api/youtube/v3, then give the youtube link to AudioFromYoutubeLink
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
