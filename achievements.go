package main

type Achievement struct {
	ID          string
	Name        string
	Description string
	Check       func(Stats) bool
}

var achievementDefs = []Achievement{
	{
		ID: "needle_drop", Name: "Needle Drop",
		Description: "Scrobbled your very first track.",
		Check:       func(s Stats) bool { return s.TotalScrobbles >= 1 },
	},
	{
		ID: "finding_the_groove", Name: "Finding the Groove",
		Description: "Reached 10 scrobbles.",
		Check:       func(s Stats) bool { return s.TotalScrobbles >= 10 },
	},
	{
		ID: "triple_digits", Name: "Triple Digits",
		Description: "Reached 100 scrobbles.",
		Check:       func(s Stats) bool { return s.TotalScrobbles >= 100 },
	},
	{
		ID: "halfway_to_obsession", Name: "Halfway to Obsession",
		Description: "Reached 500 scrobbles.",
		Check:       func(s Stats) bool { return s.TotalScrobbles >= 500 },
	},
	{
		ID: "thousand_track_stare", Name: "Thousand Track Stare",
		Description: "Reached 1,000 scrobbles.",
		Check:       func(s Stats) bool { return s.TotalScrobbles >= 1000 },
	},
	{
		ID: "endless_rotation", Name: "Endless Rotation",
		Description: "Reached 5,000 scrobbles.",
		Check:       func(s Stats) bool { return s.TotalScrobbles >= 5000 },
	},
	{
		ID: "midnight_mixtape", Name: "Midnight Mixtape",
		Description: "Scrobbled 20 tracks between midnight and 4am.",
		Check:       func(s Stats) bool { return s.NightOwlCount >= 20 },
	},
	{
		ID: "sunrise_session", Name: "Sunrise Session",
		Description: "Scrobbled 20 tracks between 5am and 7am.",
		Check:       func(s Stats) bool { return s.EarlyBirdCount >= 20 },
	},
	{
		ID: "weekend_rotation", Name: "Weekend Rotation",
		Description: "Scrobbled 50 tracks on weekends.",
		Check:       func(s Stats) bool { return s.WeekendCount >= 50 },
	},
	{
		ID: "devoted_fan", Name: "Devoted Fan",
		Description: "Played the same artist 50 times.",
		Check:       func(s Stats) bool { return maxVal(s.ArtistCounts) >= 50 },
	},
	{
		ID: "stuck_on_repeat", Name: "Stuck on Repeat",
		Description: "Played the same track 20 times.",
		Check:       func(s Stats) bool { return maxVal(s.TrackCounts) >= 20 },
	},
	{
		ID: "crate_digger", Name: "Crate Digger",
		Description: "Scrobbled tracks from 100 different artists.",
		Check:       func(s Stats) bool { return len(s.ArtistCounts) >= 100 },
	},
	{
		ID: "deep_cuts", Name: "Deep Cuts",
		Description: "Scrobbled tracks from 50 different albums.",
		Check:       func(s Stats) bool { return len(s.AlbumCounts) >= 50 },
	},
	{
		ID: "all_day_listener", Name: "All-Day Listener",
		Description: "Accumulated over 24 hours of total listening time.",
		Check:       func(s Stats) bool { return s.TotalSeconds >= 24*3600 },
	},
	{
		ID: "seven_day_spin", Name: "Seven-Day Spin",
		Description: "Listened to music 7 days in a row.",
		Check:       func(s Stats) bool { return s.LongestStreak >= 7 },
	},
	{
		ID: "habit_formed", Name: "Habit Formed",
		Description: "Listened to music 30 days in a row.",
		Check:       func(s Stats) bool { return s.LongestStreak >= 30 },
	},
	{
		ID: "extended_play", Name: "Extended Play",
		Description: "Scrobbled a track 10 minutes or longer.",
		Check:       func(s Stats) bool { return s.HasLongTrack },
	},
}


var achievementByID = func() map[string]Achievement {
	m := make(map[string]Achievement, len(achievementDefs))
	for _, a := range achievementDefs {
		m[a.ID] = a
	}
	return m
}()

func evaluate(s Stats, unlocked map[string]int64) []Achievement {
	var newly []Achievement
	for _, a := range achievementDefs {
		if _, ok := unlocked[a.ID]; ok {
			continue
		}
		if a.Check(s) {
			newly = append(newly, a)
		}
	}
	return newly
}
