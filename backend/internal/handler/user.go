package handler

import (
	"net/http"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/handler/dto"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type userHandler struct {
	userService domain.UserService
	validate    *validator.Validate
}

func NewUserHandler(r chi.Router, authMW *middleware.AuthMiddleware, userService domain.UserService) {
	h := userHandler{
		userService: userService,
		validate:    validator.New(),
	}

	// Чтение доступно любым авторизованным пользователям с правом view
	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth)
		r.Get("/api/v1/users", h.GetList)
		r.Get("/api/v1/users/{id}", h.GetByID)
	})

	// Запись доступна только администраторам
	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth, authMW.RequireAdmin)
		r.Post("/api/v1/users", h.Create)
		r.Put("/api/v1/users/{id}", h.Update)
		r.Delete("/api/v1/users/{id}", h.SoftDelete)
		r.Delete("/api/v1/users/{id}/permanent", h.HardDelete)
		r.Patch("/api/v1/users/{id}/active", h.SetActive)
	})
}

func (h *userHandler) GetList(w http.ResponseWriter, r *http.Request) {
	pagination := httputils.ParsePagination(r)

	filter := domain.UserFilter{
		FirstName:    httputils.GetQueryValue(r, "first_name"),
		SecondName:   httputils.GetQueryValue(r, "second_name"),
		LastName:     httputils.GetQueryValue(r, "last_name"),
		Email:        httputils.GetQueryValue(r, "email"),
		DepartmentID: httputils.ParseUUIDQuery(r, "department_id"),
		RoleID:       httputils.ParseUUIDQuery(r, "role_id"),
		IsActive:     httputils.ParseFilterBool(r, "is_active"),
		Pagination:   pagination,
	}

	result, err := h.userService.GetList(r.Context(), filter)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	usersDTO := dto.UsersToResponseSlice(result.Data)
	response := dto.NewPagedResponse(usersDTO, result.Pagination)
	httputils.WriteOk(w, response)
}

func (h *userHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	user, err := h.userService.GetByIDWithDetail(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.UserToResponse(user))
}

// ---------- Операции (только для админов) ----------

func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, err := httputils.DecodeJSONBodyWithValidate[dto.UserCreateRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	user := domain.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: req.Password, // в сервисе преобразуется в хеш
	}

	if req.SecondName != nil {
		user.SecondName = *req.SecondName
	}
	if req.DepartmentID != nil {
		id, err := uuid.Parse(*req.DepartmentID)
		if err == nil {
			user.DepartmentID = &id
		}
	}
	if req.RoleID != nil {
		id, err := uuid.Parse(*req.RoleID)
		if err == nil {
			user.RoleID = &id
		}
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	} else {
		user.IsActive = true // по умолчанию активен
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	if err := h.userService.Create(r.Context(), &user); err != nil {
		httputils.WriteError(w, err)
		return
	}

	// После создания возвращаем полную информацию через GetByIDWithDetail
	created, err := h.userService.GetByIDWithDetail(r.Context(), user.ID)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	httputils.WriteJSON(w, http.StatusCreated, dto.UserToResponse(created))
}

func (h *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	req, err := httputils.DecodeJSONBodyWithValidate[dto.UserUpdateRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Получаем существующего пользователя (с деталями не обязательно, достаточно базовой модели)
	existing, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Применяем только переданные поля
	if req.FirstName != nil {
		existing.FirstName = *req.FirstName
	}
	if req.SecondName != nil {
		existing.SecondName = *req.SecondName
	}
	if req.LastName != nil {
		existing.LastName = *req.LastName
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.DepartmentID != nil {
		if *req.DepartmentID == "" {
			existing.DepartmentID = nil
		} else {
			depID, err := uuid.Parse(*req.DepartmentID)
			if err == nil {
				existing.DepartmentID = &depID
			}
		}
	}
	if req.RoleID != nil {
		if *req.RoleID == "" {
			existing.RoleID = nil
		} else {
			roleID, err := uuid.Parse(*req.RoleID)
			if err == nil {
				existing.RoleID = &roleID
			}
		}
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if req.IsAdmin != nil {
		existing.IsAdmin = *req.IsAdmin
	}

	if err := h.userService.Update(r.Context(), existing); err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Возвращаем обновлённого пользователя с деталями
	updated, err := h.userService.GetByIDWithDetail(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	httputils.WriteOk(w, dto.UserToResponse(updated))
}

func (h *userHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := h.userService.SoftDelete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "пользователь помечен удалённым"})
}

func (h *userHandler) HardDelete(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := h.userService.HardDelete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "пользователь удалён"})
}

func (h *userHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	req, err := httputils.DecodeJSONBodyWithValidate[dto.UserSetActiveRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := h.userService.SetIsActive(r.Context(), id, req.IsActive); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "статус активности изменён"})
}