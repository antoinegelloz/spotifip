package supabase

type Track struct {
	Name      string   `json:"name"`
	Artists   []string `json:"artists"`
	SpotifyID string   `json:"spotify_id"`
}
