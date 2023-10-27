package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/antoinegelloz/spotifip/model/supabase"
	"os"

	"github.com/antoinegelloz/spotifip/logger"
	"github.com/antoinegelloz/spotifip/storage"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type trackStore[T any] struct {
	DB *bun.DB
}

func NewTrackStore[T any]() (storage.TrackStore[T], error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(getEnvVar("POSTGRES_ADDR")),
		pgdriver.WithUser(getEnvVar("POSTGRES_USER")),
		pgdriver.WithPassword(getEnvVar("POSTGRES_PASSWORD")),
		pgdriver.WithDatabase(getEnvVar("POSTGRES_DB")),
	)

	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New())
	if viper.GetBool("debug") {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "bun.DB.Ping")
	}

	if _, err := db.NewCreateTable().Model((*T)(nil)).
		IfNotExists().Exec(context.Background()); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return trackStore[T]{DB: db}, nil
}

func (s trackStore[T]) InsertTrack(track T) error {
	if _, err := s.DB.NewInsert().Model(&track).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func (s trackStore[T]) GetLastTrack() (*supabase.Track, error) {
	rows, err := s.DB.NewSelect().Model((*T)(nil)).
		Order("id DESC").
		Limit(1).
		Rows(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var tracks []*supabase.Track
	for rows.Next() {
		var track supabase.Track
		if err := s.DB.ScanRow(context.Background(), rows, &track); err != nil {
			return nil, errors.Wrap(err, "scan track")
		}
		tracks = append(tracks, &track)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "while reading")
	}
	if len(tracks) == 0 {
		return nil, nil
	}
	if len(tracks) > 1 {
		return nil, fmt.Errorf("%d tracks found", len(tracks))
	}

	return tracks[0], nil
}

func (s trackStore[T]) Close() {
	_ = s.DB.Close()
}

func getEnvVar(key string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		err := "couldn't find env var " + key
		logger.Get().Errorf(err)
		panic(err)
	}
	return envVar
}
