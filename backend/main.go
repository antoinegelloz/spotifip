package main

import (
	"encoding/base64"
	"fmt"
	"github.com/antoinegelloz/spotifip/logger"
	"github.com/antoinegelloz/spotifip/model/fip"
	"github.com/antoinegelloz/spotifip/model/spotify"
	"github.com/antoinegelloz/spotifip/model/supabase"
	"github.com/antoinegelloz/spotifip/storage"
	"github.com/antoinegelloz/spotifip/storage/postgres"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"os"
	"strings"
	"sync"
)

func main() {
	if err := godotenv.Load(); err != nil {
		err = fmt.Errorf("couldn't load .env file: %w", err)
		logger.Get().Errorf(err.Error())
		panic(err)
	}

	c, err := getSpotifyClientCredentials()
	if err != nil {
		logger.Get().Errorf("get Spotify client credentials: %s", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(12)
	go func() {
		defer wg.Done()
		fipFip(c)
	}()
	go func() {
		defer wg.Done()
		fipElectro(c)
	}()
	go func() {
		defer wg.Done()
		fipRock(c)
	}()
	go func() {
		defer wg.Done()
		fipReggae(c)
	}()
	go func() {
		defer wg.Done()
		fipHipHop(c)
	}()
	go func() {
		defer wg.Done()
		fipNouveautes(c)
	}()
	go func() {
		defer wg.Done()
		fipSacreFrancais(c)
	}()
	go func() {
		defer wg.Done()
		fipJazz(c)
	}()
	go func() {
		defer wg.Done()
		fipGroove(c)
	}()
	go func() {
		defer wg.Done()
		fipWorld(c)
	}()
	go func() {
		defer wg.Done()
		fipMetal(c)
	}()
	go func() {
		defer wg.Done()
		fipPop(c)
	}()
	wg.Wait()
}

func fipFip(c *spotify.ClientCredentials) {
	genre := "fip"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.Fip]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.Fip](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.Fip{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipElectro(c *spotify.ClientCredentials) {
	genre := "fip_electro"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipElectro]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipElectro](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipElectro{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipPop(c *spotify.ClientCredentials) {
	genre := "fip_pop"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipPop]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipPop](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipPop{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipMetal(c *spotify.ClientCredentials) {
	genre := "fip_metal"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipMetal]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipMetal](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipMetal{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipWorld(c *spotify.ClientCredentials) {
	genre := "fip_world"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipWorld]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipWorld](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipWorld{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipGroove(c *spotify.ClientCredentials) {
	genre := "fip_groove"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipGroove]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipGroove](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipGroove{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipJazz(c *spotify.ClientCredentials) {
	genre := "fip_jazz"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipJazz]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipJazz](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipJazz{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipSacreFrancais(c *spotify.ClientCredentials) {
	genre := "fip_sacre_francais"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipSacreFrancais]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipSacreFrancais](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipSacreFrancais{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipNouveautes(c *spotify.ClientCredentials) {
	genre := "fip_nouveautes"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipNouveautes]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipNouveautes](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipNouveautes{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipHipHop(c *spotify.ClientCredentials) {
	genre := "fip_hiphop"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipHipHop]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipHipHop](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipHipHop{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipReggae(c *spotify.ClientCredentials) {
	genre := "fip_reggae"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipReggae]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipReggae](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipReggae{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func fipRock(c *spotify.ClientCredentials) {
	genre := "fip_rock"
	f, err := getFip(getEnvVar("FIP_API_URL"), genre)
	if err != nil {
		logger.Get().Errorf("get %s: %s", genre, err)
		return
	}

	if f.Now.FirstLine == "Le direct" {
		return
	}

	q := "track:" + f.Now.FirstLine
	if f.Now.SecondLine != "" {
		q += " artist:" + f.Now.SecondLine
	}
	name := fmt.Sprintf("%s %s", f.Now.FirstLine, f.Now.SecondLine)

	ts, err := postgres.NewTrackStore[supabase.FipRock]()
	if err != nil {
		logger.Get().Errorf("NewTrackStore %s: %s", genre, err)
		return
	}
	t, err := getTrackToInsert[supabase.FipRock](c, ts, name, q, f)
	if err != nil {
		logger.Get().Errorf("getTrackToInsert %s: %s", genre, err)
		return
	} else if t != nil {
		if err := ts.InsertTrack(supabase.FipRock{Track: *t}); err != nil {
			logger.Get().Errorf("InsertTrack (found) %s: %s", genre, err)
			return
		}
		logger.Get().Infow("new track",
			"genre", genre,
			"query", q,
			"name", t.Name,
			"artists", strings.Join(t.Artists, ","),
			"id", t.SpotifyID)
	}
}

func getTrackToInsert[T any](c *spotify.ClientCredentials, ts storage.TrackStore[T], name, q string, f *fip.Fip) (*supabase.Track, error) {
	lastTrack, err := ts.GetLastTrack()
	if err != nil {
		return nil, fmt.Errorf("GET Supabase last track: %s", err)
	}

	if lastTrack != nil && lastTrack.SpotifyID == "" && lastTrack.Name == name {
		logger.Get().Infof("current track %s already inserted without ID", name)
		return nil, nil
	}

	req := resty.New().R().
		SetHeader("Authorization",
			fmt.Sprintf("%s %s", c.TokenType, c.AccessToken)).
		SetQueryParam("type", "track").
		SetQueryParam("include_external", "audio").
		SetQueryParam("q", q).
		SetQueryParam("limit", "1").
		SetResult(&spotify.Search{})

	resp, err := req.
		Get("https://api.spotify.com/v1/search")
	if err != nil {
		return nil, fmt.Errorf("GET Spotify search: %s", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf(
			"GET Spotify search: resty.Response.IsError: code %d: %+v",
			resp.StatusCode(), resp.Error())
	}

	s, ok := resp.Result().(*spotify.Search)
	if !ok {
		return nil, fmt.Errorf(
			"invalid GET Spotify search response: %+v", resp.Result())
	}

	if len(s.Tracks.Items) == 0 {
		logger.Get().Infow("GET Spotify search: no results",
			"query", q)
		if lastTrack == nil ||
			lastTrack != nil && lastTrack.Name != q {
			return &supabase.Track{
				Name: name,
				Raw:  *f,
			}, nil
		}
		return nil, nil
	}

	spotifyTrack := s.Tracks.Items[0]
	if lastTrack != nil && spotifyTrack.Name == lastTrack.Name {
		logger.Get().Infof(
			"Spotify track %s (%s) already inserted", spotifyTrack.Name, spotifyTrack.ID)
		return nil, nil
	}

	var artists []string
	for _, artist := range spotifyTrack.Artists {
		artists = append(artists, artist.Name)
	}

	return &supabase.Track{
		Name:      spotifyTrack.Name,
		SpotifyID: spotifyTrack.ID,
		Artists:   artists,
		Raw:       *f,
	}, nil
}

func getFip(fipAPIURL, genre string) (*fip.Fip, error) {
	resp, err := resty.New().R().
		SetResult(&fip.Fip{}).
		Get(fipAPIURL + genre)
	if err != nil {
		return nil, fmt.Errorf("GET: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("GET: %s", resp.String())
	}

	f, ok := resp.Result().(*fip.Fip)
	if !ok {
		return nil, fmt.Errorf("invalid response: %+v", resp.Result())
	}

	if f.Now.FirstLine == "" {
		return nil, fmt.Errorf("invalid response: empty first line: %+v", f.Now)
	}

	return f, nil
}

func getSpotifyClientCredentials() (*spotify.ClientCredentials, error) {
	spotifyClientID := getEnvVar("SPOTIFY_CLIENT_ID")
	spotifyClientSecret := getEnvVar("SPOTIFY_CLIENT_SECRET")
	resp, err := resty.New().R().
		SetHeader("Authorization",
			"Basic "+base64.StdEncoding.EncodeToString(
				[]byte(spotifyClientID+":"+spotifyClientSecret))).
		SetFormData(map[string]string{"grant_type": "client_credentials"}).
		SetResult(&spotify.ClientCredentials{}).
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		return nil, fmt.Errorf("POST Spotify token: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("POST Spotify token: %s", resp.String())
	}

	c, ok := resp.Result().(*spotify.ClientCredentials)
	if !ok {
		return nil, fmt.Errorf("invalid POST Spotify token response: %+v", resp.Result())
	}

	return c, nil
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
