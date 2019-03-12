package hldiscordopus

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// JoinUserVoiceChannel finds the channel of UserID and joins it
func JoinUserVoiceChannel(UserID string, s *discordgo.Session) (*discordgo.VoiceConnection, error) {
	for _, guild := range s.State.Guilds {
		for _, vc := range guild.VoiceStates {
			if vc.UserID == UserID {
				vs, err := s.ChannelVoiceJoin(vc.GuildID, vc.ChannelID, false, false)
				if err != nil {
					return nil, fmt.Errorf("Could not join voice channel, am i allowed to?")
				}
				return vs, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find user channel")
}
