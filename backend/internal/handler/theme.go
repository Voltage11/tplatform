package handler

import (
	"net/http"
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/handler/dto"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type themeHandlers struct {
	themeService    domain.ThemeService
	questionService domain.QuestionService
	validate        *validator.Validate
}

func NewThemeHandlers(r chi.Router, authMW *middleware.AuthMiddleware, themeService domain.ThemeService, questionService domain.QuestionService) {
	h := themeHandlers{
		themeService:    themeService,
		questionService: questionService,
		validate:        validator.New(),
	}

	r.Group(func(r chi.Router) {
		r.Use(authMW.RequireAuth) // все действия требуют авторизацию
		r.Get("/api/v1/themes", h.GetList)
		r.Post("/api/v1/themes", h.Create)
		r.Put("/api/v1/themes/{id}", h.Update)
		r.Delete("/api/v1/themes/{id}", h.Delete)

		r.Get("/api/v1/themes/{theme_id}/questions", h.GetQuestions)
		r.Post("/api/v1/themes/{theme_id}/questions", h.CreateQuestion)
		r.Put("/api/v1/themes/{theme_id}/questions/{id}", h.UpdateQuestion)
		r.Delete("/api/v1/themes/{theme_id}/questions/{id}", h.DeleteQuestion)
		r.Patch("/api/v1/themes/{theme_id}/questions/{id}", h.ChangeSortQuestion)
	})
}

func (h *themeHandlers) GetList(w http.ResponseWriter, r *http.Request) {
	pagination := httputils.ParsePagination(r)
	name := httputils.GetQueryValue(r, "name")

	var createdAtFrom, createdAtTo *time.Time
	if fromStr, ok := httputils.GetQueryValueWithExist(r, "created_at_from"); ok && fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			createdAtFrom = &t
		}
	}
	if toStr, ok := httputils.GetQueryValueWithExist(r, "created_at_to"); ok && toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			createdAtTo = &t
		}
	}

	isActive := httputils.ParseFilterBool(r, "is_active")
	createdByID := httputils.ParseUUIDQuery(r, "created_by_id")

	filter := domain.ThemeFilter{
		Name:          name,
		IsActive:      isActive,
		CreatedByID:   createdByID,
		CreatedAtFrom: createdAtFrom,
		CreatedAtTo:   createdAtTo,
		Pagination:    pagination,
	}

	themes, err := h.themeService.GetList(r.Context(), filter)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	themesDTO := dto.ThemesToResponseSlice(themes.Data)
	response := dto.NewPagedResponse(themesDTO, themes.Pagination)
	httputils.WriteOk(w, response)
}

func (h *themeHandlers) Create(w http.ResponseWriter, r *http.Request) {
	req, err := httputils.DecodeJSONBodyWithValidate[dto.ThemeCreateRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	theme := domain.Theme{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
		DateBegin:   req.DateBegin,
		DateEnd:     req.DateEnd,
		CheckPoint:  req.CheckPoint,
		ImgPath:     req.ImgPath,
	}

	if err := h.themeService.Create(r.Context(), &theme); err != nil {
		httputils.WriteError(w, err)
		return
	}

	themeWithDetail, err := h.themeService.GetByIDWithDetail(r.Context(), theme.ID)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.ThemeToResponse(themeWithDetail))
}

func (h *themeHandlers) Update(w http.ResponseWriter, r *http.Request) {
	req, err := httputils.DecodeJSONBodyWithValidate[dto.ThemeUpdateRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	theme, err := h.themeService.GetByID(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	theme.Name = req.Name
	theme.Description = req.Description
	theme.IsActive = req.IsActive
	theme.DateBegin = req.DateBegin
	theme.DateEnd = req.DateEnd
	theme.CheckPoint = req.CheckPoint
	theme.ImgPath = req.ImgPath

	if err := h.themeService.Update(r.Context(), theme); err != nil {
		httputils.WriteError(w, err)
		return
	}

	themeWithDetail, err := h.themeService.GetByIDWithDetail(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.ThemeToResponse(themeWithDetail))
}

func (h *themeHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ParseUUID(r, "id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := h.themeService.Delete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "Тема удалена"})
}

func (h *themeHandlers) GetQuestions(w http.ResponseWriter, r *http.Request) {
	themeID, err := httputils.ParseUUID(r, "theme_id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	questions, err := h.questionService.GetList(r.Context(), themeID)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string][]*domain.Question{"data": questions})
}

func (h *themeHandlers) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	themeID, err := httputils.ParseUUID(r, "theme_id")
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	questionCreateRequest, err := httputils.DecodeJSONBodyWithValidate[dto.QuestionCreateRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	// Наша внутренняя валидация
	if err := questionCreateRequest.Validate(); err != nil {
		httputils.WriteError(w, err)
		return
	}

	question := domain.Question{
		ThemeID:      themeID,
		QuestionType: dto.GetQuestionTypeFromStr(questionCreateRequest.Name),
		Name:         questionCreateRequest.Name,
		PointCorrect: questionCreateRequest.PointCorrect,
	}

	if err := h.questionService.Create(r.Context(), &question); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.QuestionToResponse(&question))
}

func (h *themeHandlers) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	themeID, err := httputils.ParseUUID(r, "theme_id")
	_ = themeID
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	id, err := httputils.ParseUUID(r, "id")

	questionUpdateRequest, err := httputils.DecodeJSONBodyWithValidate[dto.QuestionUpdateRequest](r, h.validate)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}
	// Наша внутренняя валидация
	if err := questionUpdateRequest.Validate(); err != nil {
		httputils.WriteError(w, err)
		return
	}

	question, err := h.questionService.GetByID(r.Context(), id)
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	question.QuestionType = dto.GetQuestionTypeFromStr(questionUpdateRequest.Name)
	question.Name = questionUpdateRequest.Name
	question.PointCorrect = questionUpdateRequest.PointCorrect

	if err := h.questionService.Update(r.Context(), question); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, dto.QuestionToResponse(question))
}

func (h *themeHandlers) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	themeID, err := httputils.ParseUUID(r, "theme_id")
	_ = themeID
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	id, err := httputils.ParseUUID(r, "id")

	if err := h.questionService.Delete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "Вопрос удален"})
}

func (h *themeHandlers) ChangeSortQuestion(w http.ResponseWriter, r *http.Request) {
	themeID, err := httputils.ParseUUID(r, "theme_id")
	_ = themeID
	if err != nil {
		httputils.WriteError(w, err)
		return
	}

	id, err := httputils.ParseUUID(r, "id")

	if err := h.questionService.Delete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	questionSetOrderRequest, err := httputils.DecodeJSONBodyWithValidate[dto.QuestionSetOrderRequest](r, h.validate)
	if err := h.questionService.Delete(r.Context(), id); err != nil {
		httputils.WriteError(w, err)
		return
	}

	if err := h.questionService.SetSortOrder(r.Context(), id, questionSetOrderRequest.SortOrder); err != nil {
		httputils.WriteError(w, err)
		return
	}

	httputils.WriteOk(w, map[string]string{"message": "Изменение сортировки успешно"})
}
