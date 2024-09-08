package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/guneyeroglu/cine-corn-be/database"
	"github.com/guneyeroglu/cine-corn-be/middleware"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func validateInput(user *models.User) string {
	if strings.Contains(user.Username, " ") {
		return "Username must not contain spaces"

	}
	if strings.Contains(user.Password, " ") {
		return "Password must not contain spaces"

	}

	if len(user.Username) < 3 || len(user.Username) > 20 {
		return "Username must be between 3 and 20 characters"
	}

	if len(user.Password) < 8 || len(user.Password) > 16 {
		return "Password must be between 8 and 16 characters"
	}

	return ""
}

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return utils.Response(c, nil, fiber.StatusBadRequest, "Cannot parse JSON: "+err.Error())
	}

	if err := validateInput(&user); err != "" {
		return utils.Response(c, nil, fiber.StatusBadRequest, err)
	}

	var existingUser models.User
	result := database.DB.Where("username = ?", user.Username).Find(&existingUser)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	if result.RowsAffected != 0 {
		return utils.Response(c, nil, fiber.StatusConflict, "User already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.RoleID = 2
	if err := database.DB.Create(&user).Error; err != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Failed to create user")
	}

	return utils.Response(c, nil, fiber.StatusCreated, "User created successfully")
}

func Login(c *fiber.Ctx) error {
	var user (models.User)
	if err := c.BodyParser(&user); err != nil {
		return utils.Response(c, nil, fiber.StatusBadRequest, "Cannot parse JSON: "+err.Error())
	}

	if err := validateInput(&user); err != "" {
		return utils.Response(c, nil, fiber.StatusBadRequest, err)
	}

	var existingUser models.User
	result := database.DB.Joins("Role").Where("username = ?", user.Username).Find(&existingUser)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "Username not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Incorrect password")
	}

	res, err := middleware.GenerateJwt(&existingUser)
	if err != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, err.Error())
	}

	return utils.Response(
		c,
		fiber.Map{
			"token": res,
		},
		fiber.StatusOK,
		"Successfully logged in",
	)
}

func GetAuthUser(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Token is required")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	secretKey := viper.GetString("JWT_SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return utils.Response(c, claims, fiber.StatusOK, "Authentication successful")
	}

	return utils.Response(c, nil, fiber.StatusUnauthorized, "Invalid or expired token")
}
