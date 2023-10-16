package spotify

type ClientCredentials struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIN   int    `json:"expires_in"`
}
