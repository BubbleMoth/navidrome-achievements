package main

import (
	"fmt"

	"github.com/extism/go-pdk"
	"github.com/navidrome/navidrome/plugins/pdk/go/scrobbler"
)

type achievementsPlugin struct{}
func (p *achievementsPlugin) IsAuthorized(req scrobbler.IsAuthorizedRequest) (bool, error) {
	return true, nil
}

func (p *achievementsPlugin) NowPlaying(req scrobbler.NowPlayingRequest) error {
	rememberNowPlayingUser(req.Track.ID, req.Username)
	return nil
}

func (p *achievementsPlugin) PlaybackReport(req scrobbler.PlaybackReportRequest) error {
	rememberNowPlayingUser(req.Track.ID, req.Username)
	return nil
}

func (p *achievementsPlugin) Scrobble(req scrobbler.ScrobbleRequest) error {
	username := req.Username
	track := req.Track
	if username == "" {
		if cached, ok := recoverNowPlayingUser(track.ID); ok {
			username = cached
		}
	}
	if username == "" {
		pdk.Log(pdk.LogWarn, fmt.Sprintf(
			"[Achievements] scrobble for %q has no username and no cached now-playing match, skipping",
			track.Title,
		))
		return nil
	}

	stats := loadStats(username)
	applyScrobble(&stats, track.Artist, track.Album, track.Title, float64(track.Duration), req.Timestamp)
	saveStats(username, stats)

	unlocked := loadUnlocked(username)
	newly := evaluate(stats, unlocked)
	if len(newly) == 0 {
		return nil
	}

	for _, a := range newly {
		unlocked[a.ID] = req.Timestamp
		pdk.Log(pdk.LogInfo, fmt.Sprintf(
			"[Achievement] %s unlocked \"%s\" %s (total scrobbles: %d)",
			username, a.Name, a.Description, stats.TotalScrobbles,
		))
	}
	saveUnlocked(username, unlocked)
	syncPlaylist(username, unlocked)

	return nil
}

func init() {
	scrobbler.Register(&achievementsPlugin{})
}

func main() {}
