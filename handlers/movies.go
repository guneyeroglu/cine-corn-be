package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guneyeroglu/cine-corn-be/database"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
	"github.com/spf13/viper"
)

func setMoviesDefaultState(movies []models.Movie) {
	for i := range movies {
		movies[i].IsFavorite = false
		movies[i].IsAddedToList = false
	}
}

func setMovieDefaultState(movie *models.MovieDetails) {
	movie.IsFavorite = false
	movie.IsAddedToList = false
}

func parseToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	secretKey := viper.GetString("JWT_SECRET_KEY")
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
}

func getUserMovieIDs(userID interface{}) (map[uuid.UUID]bool, map[uuid.UUID]bool, error) {
	var favoriteMovieIDs, listMovieIDs []uuid.UUID

	favErr := database.DB.Model(&models.UserFavoriteMovie{}).
		Where("user_id = ?", userID).
		Pluck("movie_id", &favoriteMovieIDs).Error

	listErr := database.DB.Model(&models.UserListMovie{}).
		Where("user_id = ?", userID).
		Pluck("movie_id", &listMovieIDs).Error

	if favErr != nil || listErr != nil {
		return nil, nil, favErr
	}

	favMap := make(map[uuid.UUID]bool)
	listMap := make(map[uuid.UUID]bool)

	for _, favID := range favoriteMovieIDs {
		favMap[favID] = true
	}
	for _, listID := range listMovieIDs {
		listMap[listID] = true
	}

	return favMap, listMap, nil
}

func GetMoviesList(c *fiber.Ctx) error {
	var movies []models.Movie

	query := database.DB
	isFeatured := c.Query("isFeatured")
	isNew := c.Query("isNew")
	if isFeatured == "true" {
		query = query.Where("is_featured = ?", true)
	}

	if isNew == "true" {
		query = query.Where("is_new = ?", true)
	}

	result := query.Find(&movies)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "An error occurred while retrieving movies.")
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No movies found.")
	}

	token, err := parseToken(c)
	if err != nil || !token.Valid {
		setMoviesDefaultState(movies)
		return utils.Response(c, movies, fiber.StatusOK, "Movies retrieved successfully.")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userId"]
		favMap, listMap, err := getUserMovieIDs(userID)
		if err != nil {
			setMoviesDefaultState(movies)
			return utils.Response(c, movies, fiber.StatusOK, "Movies retrieved successfully, but failed to retrieve favorite or list movies.")
		}

		for i := range movies {
			movies[i].IsFavorite = favMap[movies[i].ID]
			movies[i].IsAddedToList = listMap[movies[i].ID]
		}
	}

	return utils.Response(c, movies, fiber.StatusOK, "Movies retrieved successfully.")
}

func GetMovieDetailsList(c *fiber.Ctx) error {
	var movieDetails models.MovieDetails
	var movieDetailsRequest models.MovieDetailsRequest

	if err := c.BodyParser(&movieDetailsRequest); err != nil {
		return utils.Response(c, nil, fiber.StatusBadRequest, "Invalid request format.")
	}

	result := database.DB.Preload("Genres").Where("id = ?", movieDetailsRequest.ID).Find(&movieDetails)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "An error occurred while retrieving movie details.")
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "Movie details not found.")
	}

	var genreNames []string
	for _, genre := range movieDetails.Genres {
		genreNames = append(genreNames, genre.Name)
	}
	movieDetails.GenreNames = genreNames

	token, err := parseToken(c)
	if err != nil || !token.Valid {
		setMovieDefaultState(&movieDetails)
		return utils.Response(c, movieDetails, fiber.StatusOK, "Movie details retrieved successfully.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return utils.Response(c, movieDetails, fiber.StatusOK, "Movie details retrieved successfully.")
	}

	userID := claims["userId"]

	var favoriteCount, listCount int64
	favErr := database.DB.Model(&models.UserFavoriteMovie{}).
		Where("user_id = ? AND movie_id = ?", userID, movieDetailsRequest.ID).
		Count(&favoriteCount).Error

	listErr := database.DB.Model(&models.UserListMovie{}).
		Where("user_id = ? AND movie_id = ?", userID, movieDetailsRequest.ID).
		Count(&listCount).Error

	if favErr != nil || listErr != nil {
		setMovieDefaultState(&movieDetails)
		return utils.Response(c, movieDetails, fiber.StatusOK, "Movie details retrieved successfully.")
	}

	movieDetails.IsFavorite = favoriteCount > 0
	movieDetails.IsAddedToList = listCount > 0

	return utils.Response(c, movieDetails, fiber.StatusOK, "Movie details retrieved successfully.")
}
