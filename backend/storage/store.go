package storage

import "github.com/antoinegelloz/spotifip/model/supabase"

type TrackStore interface {
	InsertOneTrack(track supabase.Track) error
	GetLastTrack() (*supabase.Track, error)
	Close()
}
