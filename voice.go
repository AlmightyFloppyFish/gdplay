package hldiscordopus

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// JoinUserVoiceChannel finds the channel of UserID and joins it
func JoinUserVoiceChannel(UserID string, s *discordgo.Session) (*discordgo.VoiceState, error) {
	for _, guild := range s.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == UserID {
				return vs, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find user channel")
}
