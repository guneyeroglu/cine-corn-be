package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guneyeroglu/cine-corn-be/database"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
	"gorm.io/gorm"
)

func GetFavoriteMovies(c *fiber.Ctx) error {
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

	var favoriteMovieIDs []uuid.UUID
	result := database.DB.Model(&models.UserFavoriteMovie{}).
		Where("user_id = ?", userID).
		Pluck("movie_id", &favoriteMovieIDs)

	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Failed to retrieve favorite movies")
	}

	if len(favoriteMovieIDs) == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No favorite movies found")
	}

	result = database.DB.Where("id IN ?", favoriteMovieIDs).Find(&movies)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Failed to retrieve movies")
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No movies found")
	}

	var listMovieIDs []uuid.UUID
	result = database.DB.Model(&models.UserListMovie{}).
		Where("user_id = ?", userID).
		Pluck("movie_id", &listMovieIDs)

	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Failed to retrieve list movies")
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

	return utils.Response(c, movies, fiber.StatusOK, "Favorite movies retrieved successfully")
}

func ToggleFavoriteMovie(c *fiber.Ctx) error {
	token, err := parseToken(c)
	if err != nil || !token.Valid {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Unauthorized: Invalid or missing token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return utils.Response(c, nil, fiber.StatusUnauthorized, "Unauthorized: Invalid or expired token.")
	}

	userID := claims["userId"]
	var favoriteMovie models.UserFavoriteMovie
	var favoriteMovieRequest models.UserFavoriteMovieRequest

	if err := c.BodyParser(&favoriteMovieRequest); err != nil {
		return utils.Response(c, nil, fiber.StatusBadRequest, "Bad Request: Unable to parse request.")
	}

	err = database.DB.Where("user_id = ? AND movie_id = ?", userID, favoriteMovieRequest.MovieID).First(&favoriteMovie).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := database.DB.Create(&favoriteMovieRequest).Error; err != nil {
				return utils.Response(c, nil, fiber.StatusInternalServerError, "Server Error: Unable to add movie to favorites.")
			}

			return utils.Response(c, nil, fiber.StatusCreated, "Movie added to favorites")
		}
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Server Error: Unable to retrieve favorite movie.")
	}

	if err := database.DB.Where("user_id = ? AND movie_id = ?", userID, favoriteMovieRequest.MovieID).Delete(&models.UserFavoriteMovie{}).Error; err != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Server Error: Unable to remove movie from favorites.")
	}

	return utils.Response(c, nil, fiber.StatusOK, "Movie removed from favorites")
}
