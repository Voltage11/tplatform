package service

import (
	"context"

	"github.com/Voltage11/tplatform/internal/domain"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/google/uuid"
)

type questionService struct {
	questionRepository domain.QuestionRepository
	permissionService  domain.PermissionService
	txManager          domain.Transactor
}

func NewQuestionService(questionRepository domain.QuestionRepository, permissionService domain.PermissionService, txManager domain.Transactor) domain.QuestionService {
	return &questionService{
		questionRepository: questionRepository,
		permissionService:  permissionService,
		txManager:          txManager,
	}
}

func (q *questionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	if !q.permissionService.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	return q.questionRepository.GetByID(ctx, id)
}

func (q *questionService) Create(ctx context.Context, question *domain.Question) error {
	if !q.permissionService.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionCreate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}
	// В транзакции найдем кол. вопросов для сортировочного номера и далее вставим вопрос
	return q.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Кол. вопросов в теме TODO: хз пока не знаю, может не кол. вопросов, а макс. значение искать, после фронта проверим
		cnt, err := q.questionRepository.GetLastPositionInTheme(txCtx, question.ThemeID)
		if err != nil {
			return err
		}
		// Увеличим сортировочний индекс на единицу, пока так, дале проверим на макс значение, еще проверим смену сортировки
		question.SortOrder = cnt + 1

		return q.questionRepository.Create(txCtx, question)
	})
}

func (q *questionService) Update(ctx context.Context, question *domain.Question) error {
	if !q.permissionService.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionUpdate.Name) {
		return apperror.NewForbiddenWithoutErr()
	}

	return q.questionRepository.Update(ctx, question)
}

func (s *questionService) Delete(ctx context.Context, id uuid.UUID) error {
    if !s.permissionService.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionHardDelete.Name) {
        return apperror.NewForbiddenWithoutErr()
    }

    return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
        // Находим вопрос, чтобы узнать theme_id
        question, err := s.questionRepository.GetByID(txCtx, id)
        if err != nil {
            return err
        }

        // Удаляем вопрос
        if err := s.questionRepository.Delete(txCtx, id); err != nil {
            return err
        }

        // Получаем оставшиеся вопросы темы (отсортированные)
        remaining, err := s.questionRepository.GetList(txCtx, question.ThemeID)
        if err != nil {
            return err
        }

        // Перенумеровываем
        for i, q := range remaining {
            if err := s.questionRepository.SetSortOrder(txCtx, q.ID, i+1); err != nil {
                return err
            }
        }
        return nil
    })
}

func (q *questionService) GetList(ctx context.Context, themeID uuid.UUID) ([]*domain.Question, error) {
	if !q.permissionService.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionView.Name) {
		return nil, apperror.NewForbiddenWithoutErr()
	}

	return q.questionRepository.GetList(ctx, themeID)
}

func (s *questionService) SetSortOrder(ctx context.Context, id uuid.UUID, newPosition int) error {
    if !s.permissionService.CanFromCtx(ctx, domain.EntityThemes.Name, domain.ActionUpdate.Name) {
        return apperror.NewForbiddenWithoutErr()
    }

    return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
        // Получаем перемещаемый вопрос
        question, err := s.questionRepository.GetByID(txCtx, id)
        if err != nil {
            return err
        }
        if newPosition < 1 {
            return apperror.NewBadRequest("позиция должна быть >= 1", nil)
        }

        // Получаем все вопросы темы, уже отсортированные по sort_order
        allQuestions, err := s.questionRepository.GetList(txCtx, question.ThemeID)
        if err != nil {
            return err
        }

        // Строим новый список
        var newList []*domain.Question
        inserted := false
        for _, item := range allQuestions {
            // Пропускаем перемещаемый вопрос (он будет вставлен позже)
            if item.ID == id {
                continue
            }
            // Если мы ещё не вставили вопрос и достигли нужной позиции (позиция = индекс+1)
            if !inserted && len(newList)+1 == newPosition {
                newList = append(newList, question)
                inserted = true
            }
            newList = append(newList, item)
        }
        // Если позиция больше количества вопросов, вставляем в конец
        if !inserted {
            newList = append(newList, question)
        }

        // Перенумеровываем и сохраняем
        for i, q := range newList {
            if err := s.questionRepository.SetSortOrder(txCtx, q.ID, i+1); err != nil {
                return err
            }
        }
        return nil
    })
}
