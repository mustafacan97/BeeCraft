package handlers

import (
	"context"
	"fmt"
	"platform/internal/application/ports/repositories"
	eventBus "platform/internal/application/ports/services"
	"platform/internal/domain"
	domainEvents "platform/internal/domain/events"
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
	eventBus       eventBus.EventBus
	userRepository repositories.UserRepository
	roleRepository repositories.RoleRepository
}

func NewRegisterHandler(eventBus *eventBus.EventBus, userRepository *repositories.UserRepository, roleRepository *repositories.RoleRepository) *RegisterHandler {
	return &RegisterHandler{
		eventBus:       *eventBus,
		userRepository: *userRepository,
		roleRepository: *roleRepository,
	}
}

func (h *RegisterHandler) Handle(ctx context.Context, req *RegisterRequest) (*Response[RegisterResponse], error) {
	registeredRole, err := h.roleRepository.GetSystemRoleByName(ctx, enum.REGISTERED)
	if err != nil {
		return nil, fmt.Errorf("an error occurred on registration process: %w", err)
	}

	user, _ := domain.NewUser(req.Email, req.Password, []domain.Role{*registeredRole})

	err = h.userRepository.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("an error occurred on registration process: %w", err)
	}

	h.eventBus.Publish(ctx, *domainEvents.NewUserCreatedDomainEvent(user.Email))

	return CreatedResponseWithoutData[RegisterResponse](), nil
}
