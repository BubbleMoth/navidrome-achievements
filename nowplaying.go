package main

import (
	"github.com/navidrome/navidrome/plugins/pdk/go/host"
)

const nowPlayingTTLSeconds = 300

func nowPlayingKey(trackID string) string { return "nowplaying:" + trackID }
func rememberNowPlayingUser(trackID, username string) {
	if trackID == "" || username == "" {
		return
	}
	_ = host.KVStoreSetWithTTL(nowPlayingKey(trackID), []byte(username), nowPlayingTTLSeconds)
}

func recoverNowPlayingUser(trackID string) (string, bool) {
	if trackID == "" {
		return "", false
	}
	value, found, err := host.KVStoreGet(nowPlayingKey(trackID))
	if err != nil || !found {
		return "", false
	}
	_ = host.KVStoreDelete(nowPlayingKey(trackID))
	return string(value), true
}
