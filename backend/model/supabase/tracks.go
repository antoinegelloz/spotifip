package supabase

import (
	"github.com/uptrace/bun"
)

type Fip struct {
	bun.BaseModel `bun:"table:fip"`
	Track
}

type FipElectro struct {
	bun.BaseModel `bun:"table:fip_electro"`
	Track
}

type FipRock struct {
	bun.BaseModel `bun:"table:fip_rock"`
	Track
}

type FipJazz struct {
	bun.BaseModel `bun:"table:fip_jazz"`
	Track
}

type FipGroove struct {
	bun.BaseModel `bun:"table:fip_groove"`
	Track
}

type FipWorld struct {
	bun.BaseModel `bun:"table:fip_world"`
	Track
}

type FipReggae struct {
	bun.BaseModel `bun:"table:fip_reggae"`
	Track
}

type FipHipHop struct {
	bun.BaseModel `bun:"table:fip_hiphop"`
	Track
}

type FipNouveautes struct {
	bun.BaseModel `bun:"table:fip_nouveautes"`
	Track
}

type FipSacreFrancais struct {
	bun.BaseModel `bun:"table:fip_sacre_francais"`
	Track
}

type FipMetal struct {
	bun.BaseModel `bun:"table:fip_metal"`
	Track
}

type FipPop struct {
	bun.BaseModel `bun:"table:fip_pop"`
	Track
}

type Track struct {
	ID        int64    `bun:",pk,autoincrement" json:"id"`
	Name      string   `json:"name"`
	Artists   []string `bun:",array" json:"artists"`
	SpotifyID string   `json:"spotify_id"`
	Favorite  bool     `json:"favorite"`

	Raw any `bun:"type:jsonb" json:"raw"`
}
