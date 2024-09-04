package middleware

import (
	"fmt"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
	"github.com/spf13/viper"
)

var jwtSecretKey []byte

func GenerateJwt(user *models.User) (string, error) {
	secretKey := viper.GetString("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not found in config or environment")
	}

	jwtSecretKey = []byte(secretKey)
	claims := jwt.MapClaims{
		"userId":   user.Id,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStringToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedStringToken, nil
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
