package auth

import (
	"context"
	"fmt"
	"platform/internal/application/handlers"
	"platform/internal/application/ports/secondary"
	"platform/internal/domain"
	"platform/internal/enum"
)

type RegisterRequest struct {
	FirstName           string `json:"firstName" validate:"max=64"`
	LastName            string `json:"lastName" validate:"max=64"`
	Email               string `json:"email" validate:"required,email"`
	Password            string `json:"password" validate:"required,min=8,max=16,password"`
	SubscribeNewsletter bool   `json:"subscribeNewsletter"`
}

type RegisterResponse struct {
}

type RegisterHandler struct {
	userRepository secondary.UserRepository
	roleRepository secondary.RoleRepository
}

func NewRegisterHandler(userRepository secondary.UserRepository, roleRepository secondary.RoleRepository) *RegisterHandler {
	return &RegisterHandler{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
}

func (h *RegisterHandler) Handle(ctx context.Context, req *RegisterRequest) (*handlers.Response[RegisterResponse], error) {
	registeredRole, err := h.roleRepository.GetSystemRoleByName(ctx, enum.REGISTERED)
	if err != nil {
		return nil, fmt.Errorf("an error occurred on registration process: %w", err)
	}

	user, _ := domain.NewUser(req.Email, req.Password, []domain.Role{*registeredRole})

	err = h.userRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("an error occurred on registration process: %w", err)
	}

	return handlers.CreatedResponseWithoutData[RegisterResponse](), nil
}
