package storage

import "github.com/antoinegelloz/spotifip/model/supabase"

type TrackStore[T any] interface {
	InsertTrack(track T) error
	GetLastTrack() (*supabase.Track, error)
	Close()
}
