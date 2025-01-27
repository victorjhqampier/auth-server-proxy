package adapters

type Response struct {
	StatusCode  int    `json:"-"`
	AccessToken string `json:"access_token,omitempty"`
	Scope       string `json:"scope,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	IDToken     string `json:"id_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
}
