package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rizkisundara/project-management-api/models/dto"
	"github.com/rizkisundara/project-management-api/services"
	"github.com/rizkisundara/project-management-api/utils"
)

var validate = validator.New()

type UserController struct {
	services services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{services: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	req := new(dto.UserRegisterRequest)

	if err := ctx.BodyParser(req); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(ctx, "Validasi Gagal", err.Error())
	}

	user, err := c.services.Register(req)
	if err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	response := dto.ToUserResponse(user)
	return utils.Created(ctx, "Registrasi Berhasil", response)
}
