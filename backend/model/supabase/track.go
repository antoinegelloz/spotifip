package supabase

type Track struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Artists   []string `json:"artists"`
	SpotifyID string   `json:"spotify_id"`
	Favorite  bool     `json:"favorite"`
	Raw       any      `json:"raw"`
}
