package handler

import (
	"belajar-backend-golang/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
)

type UserHandler struct {
	Service    *service.UserService
	Validation *validator.Validate
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		Service:    svc,
		Validation: validator.New(),
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username" validate:"required,min=5,max=25"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=25"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	if err := h.Validation.Struct(input); err != nil {
		errorJSON(w, http.StatusBadRequest, map[string]any{
			"status": "error",
			"errors": FormatValidationError(err),
		})
		return
	}
	newUser, err := h.Service.Register(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Registering account successfully",
		"data":    newUser,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=25"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	if err := h.Validation.Struct(input); err != nil {
		errorJSON(w, http.StatusBadRequest, map[string]any{
			"status": "error",
			"errors": FormatValidationError(err),
		})
		return
	}
	user, token, err := h.Service.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		errorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Logging in account successfully",
		"token":   token,
		"data":    user,
	})
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		switch e.Tag() {
		case "required":
			errors[e.Field()] = e.Field() + " is required"
		case "email":
			errors[e.Field()] = "Format email is not valid"
		case "min":
			errors[e.Field()] = e.Field() + " must be at least " + e.Param() + " characters"
		case "max":
			errors[e.Field()] = e.Field() + " must be at most " + e.Param() + " characters"
		default:
			errors[e.Field()] = e.Field() + " is invalid"
		}

	}
	return errors
}
