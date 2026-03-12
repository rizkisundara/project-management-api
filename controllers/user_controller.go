package controllers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rizkisundara/project-management-api/models"
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

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Invalid ID format", err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Failed to parse request body", err.Error())
	}

	userExisting, err := c.services.GetByPublicID(publicID.String())
	if err != nil {
		return utils.NotFound(ctx, "User Not Found", err.Error())
	}

	userExisting.Name = user.Name

	if err := c.services.Update(userExisting); err != nil {
		return utils.InternalServerError(ctx, "Failed to update user", err.Error())
	}

	userResponse := dto.ToUserResponse(userExisting)
	return utils.Success(ctx, "User updated successfully", userResponse)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.services.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Failed to delete user", err.Error())
	}

	return utils.Success(ctx, "User deleted successfully", nil)
}
