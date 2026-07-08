package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/extism/go-pdk"
	"github.com/navidrome/navidrome/plugins/pdk/go/host"
)

const achievementsPlaylistName = "🏆 Achievements"

type playlistObj struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func syncPlaylist(username string, unlocked map[string]int64) {
	comment := buildComment(unlocked)

	id, err := findPlaylistID(username)
	if err != nil {
		pdk.Log(pdk.LogWarn, fmt.Sprintf("[Achievements] could not look up playlist for %s: %v", username, err))
		return
	}
	if id == "" {
		newID, createErr := createPlaylist(username)
		if createErr != nil {
			pdk.Log(pdk.LogWarn, fmt.Sprintf("[Achievements] could not create playlist for %s: %v", username, createErr))
			return
		}
		id = newID
	}
	if updateErr := updatePlaylistComment(username, id, comment); updateErr != nil {
		pdk.Log(pdk.LogWarn, fmt.Sprintf("[Achievements] could not update playlist for %s: %v", username, updateErr))
	}
}

func buildComment(unlocked map[string]int64) string {
	type entry struct {
		a  Achievement
		ts int64
	}
	entries := make([]entry, 0, len(unlocked))
	for id, ts := range unlocked {
		if a, ok := achievementByID[id]; ok {
			entries = append(entries, entry{a, ts})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ts < entries[j].ts })

	var b strings.Builder
	fmt.Fprintf(&b, "%d/%d achievements unlocked\n\n", len(entries), len(achievementDefs))
	const maxLen = 3000 // keep it well under typical text-column limits
	for _, e := range entries {
		when := time.Unix(e.ts, 0).UTC().Format("2006-01-02")
		line := fmt.Sprintf("🏆 %s %s (%s)\n", e.a.Name, e.a.Description, when)
		if b.Len()+len(line) > maxLen {
			b.WriteString("…and more!")
			break
		}
		b.WriteString(line)
	}
	return b.String()
}

func findPlaylistID(username string) (string, error) {
	raw, err := subsonicCall(fmt.Sprintf("getPlaylists?u=%s", url.QueryEscape(username)))
	if err != nil {
		return "", err
	}
	var resp struct {
		SubsonicResponse struct {
			Playlists struct {
				Playlist json.RawMessage `json:"playlist"`
			} `json:"playlists"`
		} `json:"subsonic-response"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", err
	}
	playlists, err := parsePlaylistList(resp.SubsonicResponse.Playlists.Playlist)
	if err != nil {
		return "", err
	}
	for _, p := range playlists {
		if p.Name == achievementsPlaylistName {
			return p.ID, nil
		}
	}
	return "", nil
}


func parsePlaylistList(raw json.RawMessage) ([]playlistObj, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var list []playlistObj
	if err := json.Unmarshal(raw, &list); err == nil {
		return list, nil
	}
	var single playlistObj
	if err := json.Unmarshal(raw, &single); err == nil {
		return []playlistObj{single}, nil
	}
	return nil, fmt.Errorf("unrecognized playlist list shape")
}

func createPlaylist(username string) (string, error) {
	q := fmt.Sprintf("createPlaylist?name=%s&u=%s",
		url.QueryEscape(achievementsPlaylistName), url.QueryEscape(username))
	raw, err := subsonicCall(q)
	if err != nil {
		return "", err
	}
	var resp struct {
		SubsonicResponse struct {
			Playlist playlistObj `json:"playlist"`
		} `json:"subsonic-response"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", err
	}
	if resp.SubsonicResponse.Playlist.ID == "" {
		return "", fmt.Errorf("createPlaylist returned no id")
	}
	return resp.SubsonicResponse.Playlist.ID, nil
}

func updatePlaylistComment(username, playlistID, comment string) error {
	q := fmt.Sprintf("updatePlaylist?playlistId=%s&comment=%s&u=%s",
		url.QueryEscape(playlistID), url.QueryEscape(comment), url.QueryEscape(username))
	_, err := subsonicCall(q)
	return err
}


func subsonicCall(path string) ([]byte, error) {
	resp, err := host.SubsonicAPICall(path)
	if err != nil {
		return nil, err
	}
	return []byte(resp), nil
}
