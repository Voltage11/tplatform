package dto

import (
	"time"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
)

type QuestionCreateRequest struct {
	QuestionType string `json:"question_type" validate:"required"`
	Name         string `json:"name" validate:"required,min=3"`
	PointCorrect int    `json:"point_correct"`
}

func (q *QuestionCreateRequest) Validate() error {
	if q.QuestionType != domain.QuestionTypeMulti.String() && q.QuestionType != domain.QuestionTypeSingle.String() {
		return apperror.NewBadRequest("Тип вопроса может иметь два значения: single, multi", nil)
	}

	return nil
}

func GetQuestionTypeFromStr(name string) domain.QuestionType {
	if name == "single" {
		return domain.QuestionTypeSingle
	}

	if name == "multi" {
		return domain.QuestionTypeMulti
	}

	return domain.QuestionTypeSingle
}

type QuestionUpdateRequest struct {
	QuestionType string `json:"question_type" validate:"required"`
	Name         string `json:"name" validate:"required,min=3"`
	PointCorrect int    `json:"point_correct"`
}

func (q *QuestionUpdateRequest) Validate() error {
	if q.QuestionType != domain.QuestionTypeMulti.String() || q.QuestionType != domain.QuestionTypeSingle.String() {
		return apperror.NewBadRequest("Тип вопроса может иметь два значения: single, multi", nil)
	}

	return nil
}

type QuestionSetOrderRequest struct {
	SortOrder int `json:"sort_order" validate:"required"`
}

type QuestionResponse struct {
	ID           string `json:"id"`
	ThemeID      string `json:"theme_id"`
	QuestionType string `json:"question_type"`
	Name         string `json:"name"`
	PointCorrect int    `json:"point_correct"`
	SortOrder    int    `json:"sort_order"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func QuestionToResponse(question *domain.Question) *QuestionResponse {
	return &QuestionResponse{
		ID:           question.ID.String(),
		ThemeID:      question.ThemeID.String(),
		QuestionType: string(question.QuestionType),
		Name:         question.Name,
		PointCorrect: question.PointCorrect,
		SortOrder:    question.SortOrder,
		CreatedAt:    question.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:    question.UpdatedAt.Format(time.RFC3339Nano),
	}
}

func QuestionsToResponseSlice(questions []*domain.Question) []*QuestionResponse {
	out := make([]*QuestionResponse, len(questions))

	for i, item := range questions {
		out[i] = QuestionToResponse(item)
	}

	return out
}
