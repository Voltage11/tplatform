package repository

import (
	"context"
	"time"

	"github.com/Voltage11/tplatform/internal/db"
	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
)

type questionRepository struct {
	db *db.PostgresDB
}

func NewQuestionRepository(db *db.PostgresDB) domain.QuestionRepository {
	return &questionRepository{
		db: db,
	}
}

func (q *questionRepository)GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	query := `
		SELECT id, theme_id, question_type, name, point_correct, sort_order, created_at, updated_at
		FROM questions
		where id  =$1
	`
	var question domain.Question

	executor := q.db.GetDB(ctx)

	if err := executor.QueryRow(ctx, query, id).Scan(
		&question.ID,
		&question.ThemeID,
		&question.QuestionType,
		&question.Name,
		&question.PointCorrect,
		&question.SortOrder,
		&question.CreatedAt,
		&question.UpdatedAt); err != nil {
			return nil, apperror.NewPostgresError(err)
		}

	return &question, nil
}

func (q *questionRepository) Create(ctx context.Context, question *domain.Question) error {
	query := `
		INSERT INTO questions(theme_id, question_type, name, point_correct, sort_order, created_at, updated_at)
        VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`
	now := time.Now().UTC()
	question.CreatedAt = now
	question.UpdatedAt = now

	executor := q.db.GetDB(ctx)

	if err := executor.QueryRow(ctx, query,
		question.ThemeID,
		question.QuestionType,
		question.Name,
		question.PointCorrect,
		question.SortOrder,
		question.CreatedAt,
		question.UpdatedAt).Scan(&question.ID); err != nil {
		return apperror.NewPostgresError(err)
	}

	return nil
}

func (q *questionRepository) Update(ctx context.Context, question *domain.Question) error {
	query := `
		UPDATE questions 
			SET question_type = $1, 
				name = $2,
				point_correct = $3,				
				updated_at = $4
		WHERE id = $5
	`
	question.UpdatedAt = time.Now().UTC()

	executor := q.db.GetDB(ctx)

	result, err := executor.Exec(ctx, query,
		question.QuestionType,
		question.Name,
		question.PointCorrect,
		// question.SortOrder,
		question.UpdatedAt,
		question.ID)
	if err != nil {
		return apperror.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Вопрос не найден", nil)
	}

	return nil
}

// GetLastPositionInTheme Сколько уже вопросов в теме, для получения след. номера сортировки
func (q *questionRepository) GetLastPositionInTheme(ctx context.Context, themeID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) as questions_count
		FROM questions
		WHERE theme_id = $1
	`
	executor := q.db.GetDB(ctx)

	var out int

	if err := executor.QueryRow(ctx, query, themeID).Scan(&out); err != nil {
		return 0, apperror.NewPostgresError(err)
	}

	return out, nil
}

func (q *questionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM questions WHERE id = $1`

	executor := q.db.GetDB(ctx)

	result, err := executor.Exec(ctx, query, id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Вопрос не найден", nil)
	}

	return nil
}

func (q *questionRepository) GetList(ctx context.Context, themeID uuid.UUID) ([]*domain.Question, error) {
	query := `
		SELECT id, theme_id, question_type, name, point_correct, sort_order, created_at, updated_at
		FROM questions
		WHERE theme_id  =$1
		ORDER BY sort_order
	`
	var questions []*domain.Question

	executor := q.db.GetDB(ctx)

	result, err := executor.Query(ctx, query, themeID)	
	if err != nil {
		return questions, apperror.NewPostgresError(err)
	}

	defer result.Close()

	for result.Next() {
		var question domain.Question

		if err := result.Scan(
			&question.ID,
			&question.ThemeID,
			&question.QuestionType,
			&question.Name,
			&question.PointCorrect,
			&question.SortOrder,
			&question.CreatedAt,
			&question.UpdatedAt); err != nil {
				return questions, apperror.NewPostgresError(err)
			}

		questions = append(questions, &question)
	}

	return questions, nil
}

func (q *questionRepository) SetSortOrder(ctx context.Context, id uuid.UUID, sortOrder int) error {
	query := `
		UPDATE questions 
			SET sort_order = $1,
				updated_at = $2
		WHERE id = $3
	`
	
	executor := q.db.GetDB(ctx)

	result, err := executor.Exec(ctx, query, sortOrder, time.Now().UTC(), id)
	if err != nil {
		return apperror.NewPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("Вопрос не найден", nil)
	}

	return nil
}