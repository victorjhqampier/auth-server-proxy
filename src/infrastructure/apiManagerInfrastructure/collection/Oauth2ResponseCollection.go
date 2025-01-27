package httpinfracollection

type Oauth2ResponseCollection struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
}
