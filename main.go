package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/antoinegelloz/spotifip/logger"
	"github.com/antoinegelloz/spotifip/models/fip"
	"github.com/antoinegelloz/spotifip/models/spotify"
	"github.com/antoinegelloz/spotifip/models/supabase"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const (
	ignore1 = "Le direct"
	ignore2 = "De Air à Soulwax, de Superpoze à Tosca, gardez le kick avec notre sélection électronique"
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

	if f.Now.FirstLine == ignore1 && f.Now.SecondLine == ignore2 {
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

	lastTrack, err := getSBLastTrack()
	if err != nil {
		logger.Get().Errorf("GET Supabase last track: %s", err)
		return
	}

	searchQuery := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		searchQuery += " artist:" + f.Now.SecondLine
	}

	req := httpClient.R().EnableTrace().
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
		logger.Get().Errorf("GET Spotify search IsError %d: %+v", resp.StatusCode(), resp.Error())
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
			if err = postSB(fmt.Sprintf(
				"%s %s", f.Now.FirstLine, f.Now.SecondLine), "", []string{}, false); err != nil {
				logger.Get().Errorf("POST Supabase: %s", err)
			}
		}
		return
	}

	currTrack := s.Tracks.Items[0]
	if currTrack.Name == lastTrack.Name {
		logger.Get().Infof("current track %s (%s) already inserted", currTrack.Name, currTrack.ID)
		if lastTrack.Favorite {
			if err = setSBFavorite(currTrack.ID); err != nil {
				logger.Get().Errorf("Set favorite: %s", err)
			}
		}
		return
	}

	var artists []string
	for _, artist := range currTrack.Artists {
		artists = append(artists, artist.Name)
	}

	insertedTracks, err := getSBInsertedTracks(currTrack.ID)
	if err != nil {
		logger.Get().Errorf("GET Supabase inserted tracks: %s", err)
		return
	}

	favorite := false
	for _, t := range insertedTracks {
		if t.Favorite {
			favorite = true
		}
		if favorite && !t.Favorite {
			if err = setSBFavorite(t.SpotifyID); err != nil {
				logger.Get().Errorf("Set favorite: %s", err)
			}
			break
		}
	}

	if err = postSB(currTrack.Name, currTrack.ID, artists, favorite); err != nil {
		logger.Get().Errorf("POST Supabase: %s", err)
		return
	}

	logger.Get().Infow("new track",
		"query", searchQuery,
		"name", currTrack.Name,
		"artists", strings.Join(artists, ","),
		"id", currTrack.ID)
}

func getSBLastTrack() (supabase.Track, error) {
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
		return supabase.Track{}, fmt.Errorf("IsError %d: %+v", resp.StatusCode(), resp.Error())
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

func getSBInsertedTracks(spotifyID string) ([]supabase.Track, error) {
	resp, err := httpClient.R().EnableTrace().
		SetHeaders(map[string]string{
			"apikey":        sbSpotifipServiceKey,
			"Authorization": "Bearer " + sbSpotifipServiceKey,
		}).
		SetResult(&[]supabase.Track{}).
		Get(fmt.Sprintf("%s%s?select=*&spotify_id=eq.%s",
			sbSpotifipURL+sbAPIBaseURL, sbFipElectroDB, spotifyID))
	if err != nil {
		return []supabase.Track{}, err
	}

	if resp.IsError() {
		return []supabase.Track{}, fmt.Errorf("IsError %d: %+v", resp.StatusCode(), resp.Error())
	}

	tracks, ok := resp.Result().(*[]supabase.Track)
	if !ok {
		return []supabase.Track{}, fmt.Errorf("invalid result: %+v", resp.Result())
	}

	return *tracks, nil
}

func postSB(name, spotifyID string, artists []string, favorite bool) error {
	type post struct {
		Name      string   `json:"name"`
		Artists   []string `json:"artists"`
		SpotifyID string   `json:"spotify_id"`
		Favorite  bool     `json:"favorite"`
	}
	resp, err := httpClient.R().EnableTrace().
		SetHeaders(map[string]string{
			"apikey":        sbSpotifipServiceKey,
			"Authorization": "Bearer " + sbSpotifipServiceKey,
			"Content-Type":  "application/json",
		}).
		SetBody(post{
			Name:      name,
			Artists:   artists,
			SpotifyID: spotifyID,
			Favorite:  favorite,
		}).
		Post(fmt.Sprintf("%s%s", sbSpotifipURL+sbAPIBaseURL, sbFipElectroDB))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("IsError %d: %+v", resp.StatusCode(), resp.Error())
	}

	return nil
}

func setSBFavorite(spotifyID string) error {
	type fav struct {
		Favorite bool `json:"favorite"`
	}
	resp, err := httpClient.R().EnableTrace().
		SetHeaders(map[string]string{
			"apikey":        sbSpotifipServiceKey,
			"Authorization": "Bearer " + sbSpotifipServiceKey,
			"Content-Type":  "application/json",
			"Prefer":        "return=minimal",
		}).
		SetBody(fav{
			Favorite: true,
		}).
		Patch(fmt.Sprintf("%s%s?spotify_id=eq.%s",
			sbSpotifipURL+sbAPIBaseURL, sbFipElectroDB, spotifyID))
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("IsError %d: %+v", resp.StatusCode(), resp.Error())
	}

	contentRangeSlice := resp.RawResponse.Header.Values("Content-Range")
	if len(contentRangeSlice) != 1 {
		return fmt.Errorf("too big content range slice: %+v", contentRangeSlice)
	}
	contentRange := contentRangeSlice[0]
	if len(contentRange) < 5 {
		return fmt.Errorf("too small content range: %s", contentRange)
	}

	l, err := strconv.Atoi(
		contentRange[2 : len(contentRange)-2])
	if err != nil {
		return fmt.Errorf("strconv.Atoi: %w", err)
	}

	logger.Get().Infof("%d tracks with Spotify ID %s favorited", l+1, spotifyID)
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
