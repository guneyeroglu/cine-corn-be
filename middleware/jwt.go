package middleware

import (
	"log"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guneyeroglu/cine-corn-be/utils"
	"github.com/spf13/viper"
)

var jwtSecretKey []byte

func GenerateJwt(UserId string) (string, error) {
	secretKey := viper.GetString("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("JWT_SECRET_KEY not found in config or environment")
	}

	jwtSecretKey = []byte(secretKey)
	claims := jwt.MapClaims{
		"userId": UserId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return utils.Response(c, nil, fiber.StatusUnauthorized, "Unauthorized")
}

func JwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: jwtSecretKey},
		ErrorHandler: jwtErrorHandler,
	})
}
