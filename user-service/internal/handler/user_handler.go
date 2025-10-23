package handler

import (
	"net/http"

	"user-service/internal/form"
	"user-service/internal/model"
	"user-service/internal/service"
	"user-service/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	UserService *service.UserService
}

// @Summary Register a new user
// @Description Register a new user with the given details
// @Tags Users
// @Accept json
// @Produce json
// @Param RegisterForm body form.RegistrationForm true "Register Form"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var form form.RegistrationForm

	if err := c.ShouldBindJSON(&form); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := util.FormatValidationError(validationErrors)
			util.RespondWithError(c, http.StatusBadRequest, "Validation failed", formattedErrors)
			return
		}
		util.RespondWithError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	err := h.UserService.RegisterUser(form.Username, form.Email, form.Password, form.FirstName, form.LastName)
	if err != nil {
		switch err {
		case service.ErrUsernameTaken:
			util.RespondWithError(c, http.StatusBadRequest, "Username already taken", nil)
			return
		case service.ErrEmailRegistered:
			util.RespondWithError(c, http.StatusBadRequest, "Email already registered", nil)
			return
		default:
			util.RespondWithError(c, http.StatusInternalServerError, "Registration failed", nil)
			return
		}
	}
	util.RespondWithSuccess(c, "Registration successful", nil)
}

// @Summary Log in a user
// @Description Log in a user with the given credentials
// @Tags Users
// @Accept json
// @Produce json
// @Param LoginForm body form.LoginForm true "Login Form"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var form form.LoginForm

	if err := c.ShouldBindJSON(&form); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := util.FormatValidationError(validationErrors)
			util.RespondWithError(c, http.StatusBadRequest, "Validation failed", formattedErrors)
			return
		}
		util.RespondWithError(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	user, err := h.UserService.LoginUser(form.Username, form.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			util.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials", nil)
			return
		default:
			validationErrors := util.FormatValidationError(err)
			util.RespondWithError(c, http.StatusBadRequest, "Validation failed", validationErrors)
			return
		}
	}

	token, err := util.GenerateJWT(user.Username)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "could not generate token", nil)
		return
	}
	util.RespondWithSuccess(c, "login successful", gin.H{"token": token})
}

// @Summary Get user profile
// @Description Get the profile of the logged-in user
// @Tags Users
// @Produce json
// @Success 200 {object} model.APIResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/profile [get]
// @Security ApiKeyAuth
func (h *UserHandler) Profile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		util.RespondWithError(c, http.StatusUnauthorized, "no username found", nil)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		util.RespondWithError(c, http.StatusUnauthorized, "Invalid username format", nil)
		return
	}

	user, err := h.UserService.Repo.GetUserByUsername(usernameStr)
	if err != nil {
		util.RespondWithError(c, http.StatusUnauthorized, "no username found", nil)
		return
	}
	util.RespondWithSuccess(c, "user data", gin.H{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"username":  user.Username,
		"email":     user.Email,
	})
}

// @Summary Update user profile
// @Description Update the profile details of an authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Param UpdateUserForm body form.UpdateUserForm true "Update User Form"
// @Success 200 {object} model.APIResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /auth/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		util.RespondWithError(c, http.StatusUnauthorized, "no username found", nil)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		util.RespondWithError(c, http.StatusUnauthorized, "Invalid username format", nil)
		return
	}
	var form form.UpdateUserForm
	if err := c.ShouldBindJSON(&form); err != nil {
		validationErrors := util.FormatValidationError(err)
		util.RespondWithError(c, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	updatedUser := model.User{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
	}

	err := h.UserService.Repo.UpdateUser(usernameStr, &updatedUser)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "could not update user profile", nil)
		return
	}
	util.RespondWithSuccess(c, "Profile updated successfully", nil)
}
