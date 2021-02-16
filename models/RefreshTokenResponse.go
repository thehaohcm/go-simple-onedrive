package models

type RefreshTokenResponse struct {
	TokenType    string `json:"token_type,omitempty"`
	Score        string `json:"scope,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	ExtExpiresIn int    `json:"ext_expires_in,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IdToken      string `json:"id_token,omitempty"`
}
