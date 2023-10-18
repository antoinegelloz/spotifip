package main

import (
	"encoding/base64"
	"fmt"
	"github.com/antoinegelloz/spotifip/storage/postgres"
	"os"
	"strings"

	"github.com/antoinegelloz/spotifip/logger"
	"github.com/antoinegelloz/spotifip/model/fip"
	"github.com/antoinegelloz/spotifip/model/spotify"
	"github.com/antoinegelloz/spotifip/model/supabase"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		err = fmt.Errorf("couldn't load .env file: %w", err)
		logger.Get().Errorf(err.Error())
		panic(err)
	}

	fipElectroURL := getEnvVar("FIP_ELECTRO_URL")
	resp, err := resty.New().R().
		SetResult(&fip.Fip{}).
		Get(fipElectroURL)
	if err != nil {
		logger.Get().Errorf("GET Fip Electro: %s", err)
		return
	}

	if resp.IsError() {
		logger.Get().Errorf("GET Fip Electro: %s", resp.String())
		return
	}

	f, ok := resp.Result().(*fip.Fip)
	if !ok {
		logger.Get().Errorf("invalid GET Fip Electro response: %+v", resp.Result())
		return
	}

	if f.Now.FirstLine == "" {
		logger.Get().Errorf("invalid GET Fip Electro response: empty first line: %+v", f.Now)
		return
	}

	if f.Now.FirstLine == "Le direct" &&
		f.Now.SecondLine == "De Air à Soulwax, de Superpoze à Tosca, gardez le kick avec notre sélection électronique" {
		return
	}

	spotifyClientID := getEnvVar("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := getEnvVar("SPOTIFY_CLIENT_SECRET")
	resp, err = resty.New().R().
		SetHeader("Authorization",
			"Basic "+base64.StdEncoding.EncodeToString(
				[]byte(spotifyClientID+":"+spotifyClientSecret))).
		SetFormData(map[string]string{"grant_type": "client_credentials"}).
		SetResult(&spotify.ClientCredentials{}).
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		logger.Get().Errorf("POST Spotify token: %s", err)
		return
	}

	if resp.IsError() {
		logger.Get().Errorf("POST Spotify token: %s", resp.String())
		return
	}

	c, ok := resp.Result().(*spotify.ClientCredentials)
	if !ok {
		logger.Get().Errorf("invalid POST Spotify token response: %+v", resp.Result())
		return
	}

	ts, err := postgres.NewTrackStore()
	if err != nil {
		logger.Get().Errorf("NewTrackStore: %s", err)
		return
	}

	lastTrack, err := ts.GetLastTrack()
	if err != nil {
		logger.Get().Errorf("GET Supabase last track: %s", err)
		return
	}

	searchQuery := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		searchQuery += " artist:" + f.Now.SecondLine
	}

	nowName := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)
	if lastTrack.SpotifyID == "" && lastTrack.Name == nowName {
		logger.Get().Infof("current track %s already inserted without ID", nowName)
		return
	}

	req := resty.New().R().
		SetHeader("Authorization",
			fmt.Sprintf("%s %s", c.TokenType, c.AccessToken)).
		SetQueryParam("type", "track").
		SetQueryParam("include_external", "audio").
		SetQueryParam("q", searchQuery).
		SetQueryParam("limit", "1").
		SetResult(&spotify.Search{})

	resp, err = req.
		Get("https://api.spotify.com/v1/search")
	if err != nil {
		logger.Get().Errorf("GET Spotify search: %s", err)
		return
	}

	if resp.IsError() {
		logger.Get().Errorf("GET Spotify search: resty.Response.IsError: code %d: %+v", resp.StatusCode(), resp.Error())
		return
	}

	s, ok := resp.Result().(*spotify.Search)
	if !ok {
		logger.Get().Errorf("invalid GET Spotify search response: %+v", resp.Result())
		return
	}

	if len(s.Tracks.Items) == 0 {
		logger.Get().Infow("GET Spotify search: no results",
			"query", searchQuery)
		if lastTrack.Name != searchQuery {
			if err := ts.InsertOneTrack(supabase.Track{
				Name: nowName,
				Raw:  *f,
			}); err != nil {
				logger.Get().Errorf("InsertOneTrack: %s", err)
			}
		}
		return
	}

	spotifyTrack := s.Tracks.Items[0]
	if spotifyTrack.Name == lastTrack.Name {
		logger.Get().Infof("Spotify track %s (%s) already inserted", spotifyTrack.Name, spotifyTrack.ID)
		return
	}

	var artists []string
	for _, artist := range spotifyTrack.Artists {
		artists = append(artists, artist.Name)
	}

	if err := ts.InsertOneTrack(supabase.Track{
		Name:      spotifyTrack.Name,
		SpotifyID: spotifyTrack.ID,
		Artists:   artists,
		Raw:       *f,
	}); err != nil {
		logger.Get().Errorf("InsertOneTrack (found): %s", err)
		return
	}

	logger.Get().Infow("new track",
		"query", searchQuery,
		"name", spotifyTrack.Name,
		"artists", strings.Join(artists, ","),
		"id", spotifyTrack.ID)
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
