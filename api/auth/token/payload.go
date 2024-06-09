package token

import (
	"time"
)

type TokenPayload struct {
	Aud string `json:"aud"`
	Sub string `json:"sub"`
	Iss string `json:"iss"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
	Prm string `json:"prm"`
}

func NewTokenPayload(issuer, audience, subject, param string, validity int64) *TokenPayload {
	iat := time.Now().Unix()
	exp := iat + validity
	return &TokenPayload{
		Aud: audience,
		Sub: subject,
		Iss: issuer,
		Iat: iat,
		Exp: exp,
		Prm: param,
	}
}
