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
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(ctx, "Validation failed", err.Error())
	}

	user, err := c.services.Register(req)
	if err != nil {
		return utils.BadRequest(ctx, "Registration Failed", err.Error())
	}

	response := dto.ToUserResponse(user)
	return utils.Created(ctx, "Registration Successful", response)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	req := new(dto.UserLoginRequest)

	if err := ctx.BodyParser(req); err != nil {
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	if err := validate.Struct(req); err != nil {
		return utils.BadRequest(ctx, "Validation failed", err.Error())
	}

	user, err := c.services.Login(req.Email, req.Password)
	if err != nil {
		return utils.BadRequest(ctx, "Login Failed", err.Error())
	}

	token, _ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken, _ := utils.GenerateRefreshToken(user.InternalID)

	return utils.Success(ctx, "Login Successful", fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          dto.ToUserResponse(user),
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.services.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "User Not Found", err.Error())
	}

	return utils.Success(ctx, "User Found", dto.ToUserResponse(user))
}

func (c *UserController) GetAllUsersWithPagination(ctx *fiber.Ctx) error {
	query := utils.ParsePaginationQuery(ctx)

	users, total, err := c.services.GetAllUsersWithPagination(query)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to retrieve users", err.Error())
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, dto.ToUserResponse(&user))
	}

	meta := utils.BuildPaginationMeta(query, total)

	return utils.SuccessWithPagination(ctx, "Users retrieved successfully", responses, meta)
}
