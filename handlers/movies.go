package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/cine-corn-be/database"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
)

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
		return utils.Response(c, nil, fiber.StatusInternalServerError, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No movies found")
	}

	return utils.Response(c, movies, fiber.StatusOK, "Movies retrieved successfully")
}

func GetMovieDetailsList(c *fiber.Ctx) error {
	var movieDetails models.MovieDetails
	var MovieDetailsRequest models.MovieDetailsRequest
	if err := c.BodyParser(&MovieDetailsRequest); err != nil {
		return utils.Response(c, nil, fiber.StatusBadRequest, "Cannot parse JSON: "+err.Error())
	}

	result := database.DB.Preload("Genres").Where("id = ?", MovieDetailsRequest.ID).Find(&movieDetails)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No movie details found")
	}

	var genreNames []string
	for _, genre := range movieDetails.Genres {
		genreNames = append(genreNames, genre.Name)
	}

	movieDetails.GenreNames = genreNames

	return utils.Response(c, movieDetails, fiber.StatusOK, "Movie details retrieved successfully")
}
