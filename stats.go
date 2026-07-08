package main

import (
	"encoding/json"
	"strings"
	"time"
	"github.com/navidrome/navidrome/plugins/pdk/go/host"
)
type Stats struct {
	TotalScrobbles int            `json:"totalScrobbles"`
	TotalSeconds   float64        `json:"totalSeconds"`
	ArtistCounts   map[string]int `json:"artistCounts"`
	AlbumCounts    map[string]int `json:"albumCounts"`
	TrackCounts    map[string]int `json:"trackCounts"`

	LastPlayDate  string `json:"lastPlayDate"` // YYYY-MM-DD (UTC)
	CurrentStreak int    `json:"currentStreak"`
	LongestStreak int    `json:"longestStreak"`

	NightOwlCount  int  `json:"nightOwlCount"`  // scrobbles between 00:00-04:00
	EarlyBirdCount int  `json:"earlyBirdCount"` // scrobbles between 05:00-07:00
	WeekendCount   int  `json:"weekendCount"`   // scrobbles on Sat/Sun
	HasLongTrack   bool `json:"hasLongTrack"`   // ever scrobbled a track >= 10 minutes
}

func newStats() Stats {
	return Stats{
		ArtistCounts: map[string]int{},
		AlbumCounts:  map[string]int{},
		TrackCounts:  map[string]int{},
	}
}

func statsKey(username string) string { return "stats:" + username }
func unlockedKey(username string) string { return "unlocked:" + username }

func loadStats(username string) Stats {
	value, found, err := host.KVStoreGet(statsKey(username))
	if err != nil || !found {
		return newStats()
	}
	s := newStats()
	if jsonErr := json.Unmarshal(value, &s); jsonErr != nil {
		return newStats()
	}
	if s.ArtistCounts == nil {
		s.ArtistCounts = map[string]int{}
	}
	if s.AlbumCounts == nil {
		s.AlbumCounts = map[string]int{}
	}
	if s.TrackCounts == nil {
		s.TrackCounts = map[string]int{}
	}
	return s
}

func saveStats(username string, s Stats) {
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	_ = host.KVStoreSet(statsKey(username), data)
}

// loadUnlocked returns a map of achievement ID -> unix timestamp unlocked.
func loadUnlocked(username string) map[string]int64 {
	value, found, err := host.KVStoreGet(unlockedKey(username))
	m := map[string]int64{}
	if err != nil || !found {
		return m
	}
	_ = json.Unmarshal(value, &m)
	return m
}

func saveUnlocked(username string, m map[string]int64) {
	data, err := json.Marshal(m)
	if err != nil {
		return
	}
	_ = host.KVStoreSet(unlockedKey(username), data)
}

func applyScrobble(s *Stats, artist, album, track string, durationSeconds float64, timestamp int64) {
	s.TotalScrobbles++
	s.TotalSeconds += durationSeconds

	incr(s.ArtistCounts, normalize(artist))
	incr(s.AlbumCounts, normalize(album))
	incr(s.TrackCounts, normalize(artist)+" - "+normalize(track))

	if durationSeconds >= 600 {
		s.HasLongTrack = true
	}

	ts := time.Unix(timestamp, 0).UTC()
	hour := ts.Hour()
	if hour >= 0 && hour < 4 {
		s.NightOwlCount++
	}
	if hour >= 5 && hour < 7 {
		s.EarlyBirdCount++
	}
	switch ts.Weekday() {
	case time.Saturday, time.Sunday:
		s.WeekendCount++
	}

	updateStreak(s, ts.Format("2006-01-02"))
}

func updateStreak(s *Stats, dateStr string) {
	switch {
	case s.LastPlayDate == "":
		s.CurrentStreak = 1
	case dateStr == s.LastPlayDate:
		// same day, streak unchanged
	default:
		prev, errPrev := time.Parse("2006-01-02", s.LastPlayDate)
		cur, errCur := time.Parse("2006-01-02", dateStr)
		if errPrev == nil && errCur == nil && cur.After(prev) {
			if cur.Sub(prev) <= 36*time.Hour {
				s.CurrentStreak++
			} else {
				s.CurrentStreak = 1
			}
		} else {
			s.CurrentStreak = 1
		}
	}
	if s.CurrentStreak > s.LongestStreak {
		s.LongestStreak = s.CurrentStreak
	}
	s.LastPlayDate = dateStr
}

func incr(m map[string]int, key string) {
	if key == "" {
		return
	}
	m[key]++
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func maxVal(m map[string]int) int {
	max := 0
	for _, v := range m {
		if v > max {
			max = v
		}
	}
	return max
}
