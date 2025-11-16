package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"order-service/src/internal/model"
	"order-service/src/pkg/token"
	"order-service/src/pkg/utils"

	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
)

type cfgJwt struct {
	Audience  string
	Issuer    string
	Algorithm string
}

func decodeKey(secret string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func VerifyBearer(viper *viper.Viper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization", "")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return utils.Response(nil, "Invalid token!", http.StatusUnauthorized, c)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if len(tokenString) == 0 {
			return utils.Response(nil, "Invalid token!", http.StatusUnauthorized, c)
		}

		publicKey, err := decodeKey(viper.GetString("jwt.publickey"))
		if err != nil {
			return utils.Response(nil, utils.ConvertString(err), http.StatusUnauthorized, c)
		}

		parsedToken := <-token.Validate(c.Context(), publicKey, tokenString, viper.GetString("jwt.audience"), viper.GetString("jwt.issuer"), viper.GetString("jwt.signing_algorithm"))
		if parsedToken.Error != nil {
			return utils.Response(nil, utils.ConvertString(parsedToken.Error), http.StatusUnauthorized, c)
		}
		data, _ := json.Marshal(parsedToken.Data)
		var claim token.Claim
		if err := json.Unmarshal(data, &claim); err != nil {
			return utils.Response(nil, fmt.Sprintf("Invalid claim: %v", err), http.StatusUnauthorized, c)
		}
		auth := model.Auth{
			UserID:   claim.Metadata.UserID,
			FullName: claim.Metadata.FullName,
		}
		c.Locals("metadata", &auth)
		return c.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	meta := ctx.Locals("metadata")

	if meta == nil {
		return nil
	}

	auth, ok := meta.(*model.Auth)
	if !ok {
		return nil
	}

	return auth
}
