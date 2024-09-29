package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guneyeroglu/cine-corn-be/database"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
)

func GetListMovies(c *fiber.Ctx) error {
	var movies []models.Movie

	token, err := parseToken(c)
	if err != nil || !token.Valid {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Invalid or missing token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Invalid or expired token")
	}

	userID := claims["userId"]

	var listMovieIDs []uuid.UUID
	result := database.DB.Model(&models.UserListMovie{}).
		Where("user_id = ?", userID).
		Pluck("movie_id", &listMovieIDs)

	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, result.Error.Error())
	}

	if len(listMovieIDs) == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No list movies found")
	}

	result = database.DB.Where("id IN ?", listMovieIDs).Find(&movies)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No movies found")
	}

	var favoriteMovieIDs []uuid.UUID
	result = database.DB.Model(&models.UserFavoriteMovie{}).
		Where("user_id = ?", userID).
		Pluck("movie_id", &favoriteMovieIDs)

	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, result.Error.Error())
	}

	favMap := make(map[uuid.UUID]bool)
	listMap := make(map[uuid.UUID]bool)

	for _, favID := range favoriteMovieIDs {
		favMap[favID] = true
	}
	for _, listID := range listMovieIDs {
		listMap[listID] = true
	}

	for i := range movies {
		movies[i].IsFavorite = favMap[movies[i].ID]
		movies[i].IsAddedToList = listMap[movies[i].ID]
	}

	return utils.Response(c, movies, fiber.StatusOK, "List movies retrieved successfully")
}
