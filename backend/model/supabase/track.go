package supabase

import (
	"github.com/uptrace/bun"
)

type Track struct {
	bun.BaseModel `bun:"table:fip_electro"`

	ID        int64    `bun:",pk,autoincrement" json:"id"`
	Name      string   `json:"name"`
	Artists   []string `bun:",array" json:"artists"`
	SpotifyID string   `json:"spotify_id"`
	Favorite  bool     `json:"favorite"`

	Raw any `bun:"type:jsonb" json:"raw"`
}
