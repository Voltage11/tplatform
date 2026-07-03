package handler

import (
	"net/http"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/handler/dto"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type departmentHandler struct {
	departmentService domain.DepartmentService
	validate          *validator.Validate
}

func NewDepartmentHandler(r chi.Router, authMW *middleware.AuthMiddleware, departmentService domain.DepartmentService) {
	h := departmentHandler{
		departmentService: departmentService,
		validate:          validator.New(),
	}

	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth)
		r.Get("/api/v1/departments", h.GetList)
		r.Post("/api/v1/departments", h.Create)
		r.Get("/api/v1/departments/{id}", h.GetByID)
		r.Put("/api/v1/departments/{id}", h.Update)
		r.Delete("/api/v1/departments/{id}", h.SoftDelete)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth, authMW.RequireAdmin)
		r.Delete("/api/v1/departments/{id}/hard", h.HardDelete)
	})
}

func (d *departmentHandler) GetList(w http.ResponseWriter, r *http.Request) {
	paginationRequest := httputils.ParsePagination(r)

	// Получаем значение и проверяем, что оно не пустое
	nameFilter, ok := httputils.GetQueryValue(r, "name")
	if ok && nameFilter == "" {
		nameFilter = ""
	} else if !ok {
		nameFilter = ""
	}

	departments, err := d.departmentService.GetList(r.Context(), domain.DepartmentFilter{
		Name:       nameFilter,
		Pagination: paginationRequest,
	})
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	// Преобразуем домен в DTO
	deptDTOs := dto.DepartmentsToResponseSlice(departments.Data)
	response := dto.NewPagedResponse(deptDTOs, departments.Pagination)

	httputils.WriteOk(w, response)
}

func (d *departmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	department, err := d.departmentService.GetByID(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.DepartmentToResponse(department))
}

func (d *departmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	departmentCreateRequest, err := httputils.DecodeJSONBodyWithValidate[dto.DepartmentCreateRequest](r, d.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	department := domain.Department{
		Name: departmentCreateRequest.Name,
	}

	if err := d.departmentService.Create(r.Context(), &department); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.DepartmentToResponse(&department))
}

func (d *departmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	departmentUpdateRequest, err := httputils.DecodeJSONBodyWithValidate[dto.DepartmentUpdateRequest](r, d.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	department, err := d.departmentService.GetByID(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	department.Name = departmentUpdateRequest.Name

	if err := d.departmentService.Update(r.Context(), department); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.DepartmentToResponse(department))
}

func (d *departmentHandler) HardDelete(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := d.departmentService.HardDelete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "запись полностью удалена"})
}

func (d *departmentHandler) SoftDelete(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := d.departmentService.SoftDelete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "запись помечена как удалённая"})
}
