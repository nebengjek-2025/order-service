package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order-service/src/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type cfgJwt struct {
	Audience  string
	Issuer    string
	Algorithm string
}

func Validate(ctx context.Context, publicKey string, tokenString string, Audience string, Issuer string, Algorithm string) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			output <- utils.Result{Error: "Failed to parse public key"}
			return
		}
		audience := Audience
		issuer := Issuer
		algorithm := Algorithm

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{
			Audience: jwt.ClaimStrings{audience},
			Issuer:   issuer,
		}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok || token.Header["alg"] != algorithm {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return publicKey, nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				output <- utils.Result{Error: "token has been expired"}
				return
			}
			output <- utils.Result{Error: "token parsing error"}
			return
		}
		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			var tokenClaim Claim
			jsonString, _ := json.Marshal(claims)
			json.Unmarshal(jsonString, &tokenClaim)
			output <- utils.Result{Data: tokenClaim}
		} else {
			output <- utils.Result{Error: "Token is not valid!"}
			return
		}
	}()

	return output
}
