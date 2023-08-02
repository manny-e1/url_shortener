package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/manny-e1/url_shortener/model"
	"github.com/manny-e1/url_shortener/utils"
	"strconv"
)

func getAllGolies(ctx *fiber.Ctx) error {
	golies, err := model.GetAllGolies()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error getting all goly links " + err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(golies)
}

func getGoly(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not parse id " + err.Error(),
		})
	}
	goly, err := model.GetGoly(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not retrieve goly from db " + err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(goly)

}

func createGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	var goly model.Goly

	err := ctx.BodyParser(&goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}
	if goly.Random {
		goly.Goly = utils.RandomURL(8)
	}
	fmt.Println(goly)
	err = model.CreateGoly(goly)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create goly in db " + err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(goly)
}

func updateGoly(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	var goly model.Goly

	err := ctx.BodyParser(&goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}

	err = model.UpdateGoly(goly)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not update goly in db " + err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(goly)
}

func deleteGoly(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not parse id " + err.Error(),
		})
	}
	err = model.DeleteGoly(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not delete goly from db " + err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "goly deleted",
	})
}

func redirect(ctx *fiber.Ctx) error {
	golyUrl := ctx.Params("redirect")
	goly, err := model.FindByGolyUrl(golyUrl)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not find goly in DB " + err.Error(),
		})
	}
	goly.Clicked += 1
	err = model.UpdateGoly(goly)
	if err != nil {
		fmt.Printf("error updating: %v\n", err)
	}
	return ctx.Redirect(goly.Redirect, fiber.StatusTemporaryRedirect)
}

func SetupAndListen() {
	router := fiber.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type Accept",
	}))
	router.Get("/r/:redirect", redirect)
	router.Get("/goly", getAllGolies)
	router.Get("/goly/:id", getGoly)
	router.Post("/goly", createGoly)
	router.Put("/goly", updateGoly)
	router.Delete("/goly/:id", deleteGoly)
	fmt.Println("listening on 3000")
	router.Listen(":3000")
}
