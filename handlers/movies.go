package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/cine-corn-be/database"
	"github.com/guneyeroglu/cine-corn-be/models"
	"github.com/guneyeroglu/cine-corn-be/utils"
)

func GetMoviesList(c *fiber.Ctx) error {
	movies := new([]models.Movie)

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
	movieDetails := new(models.MovieDetails)
	MovieDetailsRequest := new(models.MovieDetailsRequest)
	if err := c.BodyParser(MovieDetailsRequest); err != nil {
		log.Println(err)
		return utils.Response(c, nil, fiber.StatusBadRequest, "Cannot parse JSON: "+err.Error())
	}

	result := database.DB.Where("id = ?", MovieDetailsRequest.ID).Find(&movieDetails)
	if result.Error != nil {
		return utils.Response(c, nil, fiber.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return utils.Response(c, nil, fiber.StatusNotFound, "No movie details found")
	}

	return utils.Response(c, movieDetails, fiber.StatusOK, "Movie details retrieved successfully")
}
