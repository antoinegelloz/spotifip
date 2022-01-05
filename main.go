package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/antoinegelloz/spotifip/logger"
	"github.com/antoinegelloz/spotifip/models/fip"
	"github.com/antoinegelloz/spotifip/models/spotify"
	"github.com/antoinegelloz/spotifip/models/supabase"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

var (
	httpClient           *resty.Client
	spotifyClientID      string
	spotifyClientSecret  string
	fipElectroURL        string
	sbSpotifipURL        string
	sbSpotifipServiceKey string
	sbAPIBaseURL         string
	sbFipElectroDB       string
)

func main() {
	if err := godotenv.Load(); err != nil {
		err = fmt.Errorf("couldn't load .env file: %w", err)
		logger.Get().Errorf(err.Error())
		panic(err)
	}

	httpClient = resty.New()
	spotifyClientID = getEnvVar("SPOTIFY_CLIENT_ID")
	spotifyClientSecret = getEnvVar("SPOTIFY_CLIENT_SECRET")
	fipElectroURL = getEnvVar("FIP_ELECTRO_URL")
	sbSpotifipURL = getEnvVar("SUPABASE_SPOTIFIP_URL")
	sbSpotifipServiceKey = getEnvVar("SUPABASE_SPOTIFIP_SERVICE_KEY")
	sbAPIBaseURL = getEnvVar("SUPABASE_API_BASE_URL")
	sbFipElectroDB = getEnvVar("SUPABASE_FIP_ELECTRO_DB")

	resp, err := httpClient.R().EnableTrace().
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

	resp, err = httpClient.R().EnableTrace().
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

	lastTrack, err := getSupabaseLastTrack()
	if err != nil {
		logger.Get().Errorf("GET Supabase last track: %s", err)
		return
	}

	searchQuery := f.Now.FirstLine
	if f.Now.SecondLine != "" {
		searchQuery += " " + f.Now.SecondLine
	}

	resp, err = httpClient.R().EnableTrace().
		SetHeader("Authorization",
			fmt.Sprintf("%s %s", c.TokenType, c.AccessToken)).
		SetQueryParam("type", "track").
		SetQueryParam("include_external", "audio").
		SetQueryParam("q", searchQuery).
		SetResult(&spotify.Search{}).
		Get("https://api.spotify.com/v1/search")
	if err != nil {
		logger.Get().Errorf("GET Spotify search: %s", err)
		return
	}

	if resp.IsError() {
		logger.Get().Errorf("GET Spotify search: %s", resp.String())
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
			if err = postSupabase(searchQuery, "", []string{}); err != nil {
				logger.Get().Errorf("POST Supabase: %s", err)
			}
		}
		return
	}

	currTrack := s.Tracks.Items[0]
	if currTrack.Name == lastTrack.Name {
		logger.Get().Infof("current track %s (%s) already inserted", currTrack.Name, currTrack.ID)
		return
	}

	var artists []string
	for _, artist := range currTrack.Artists {
		artists = append(artists, artist.Name)
	}

	if err = postSupabase(currTrack.Name, currTrack.ID, artists); err != nil {
		logger.Get().Errorf("POST Supabase: %s", err)
		return
	}

	logger.Get().Infow("new track",
		"query", searchQuery,
		"name", currTrack.Name,
		"artists", strings.Join(artists, ","),
		"id", currTrack.ID)
}

func getSupabaseLastTrack() (supabase.Track, error) {
	resp, err := httpClient.R().EnableTrace().
		SetHeaders(map[string]string{
			"apikey":        sbSpotifipServiceKey,
			"Authorization": "Bearer " + sbSpotifipServiceKey,
			"Range-Unit":    "items",
			"Range":         "0-0",
		}).
		SetResult(&[]supabase.Track{}).
		Get(fmt.Sprintf("%s%s?select=*&order=id.desc",
			sbSpotifipURL+sbAPIBaseURL, sbFipElectroDB))
	if err != nil {
		return supabase.Track{}, err
	}

	if resp.IsError() {
		return supabase.Track{}, fmt.Errorf("response: %s", resp.String())
	}

	lastTrack, ok := resp.Result().(*[]supabase.Track)
	if !ok {
		return supabase.Track{}, fmt.Errorf("invalid result: %+v", resp.Result())
	}

	if len(*lastTrack) == 0 {
		return supabase.Track{}, nil
	}

	if len(*lastTrack) > 1 {
		return supabase.Track{}, fmt.Errorf("more than one track: %+v", lastTrack)
	}

	return (*lastTrack)[0], nil
}

func postSupabase(name, id string, artists []string) error {
	resp, err := httpClient.R().EnableTrace().
		SetHeaders(map[string]string{
			"apikey":        sbSpotifipServiceKey,
			"Authorization": "Bearer " + sbSpotifipServiceKey,
			"Content-Type":  "application/json",
		}).
		SetBody(supabase.Track{
			Name:      name,
			Artists:   artists,
			SpotifyID: id,
		}).
		Post(fmt.Sprintf("%s%s", sbSpotifipURL+sbAPIBaseURL, sbFipElectroDB))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("response: %s", resp.String())
	}

	return nil
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
