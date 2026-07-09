package dto

import (
    "time"

    "github.com/Voltage11/tplatform/internal/domain"
)

type ThemeCreateRequest struct {
    Name        string     `json:"name" validate:"required,min=3"`
    Description string     `json:"description,omitempty"`
    IsActive    bool       `json:"is_active"`
    DateBegin   *time.Time `json:"date_begin,omitempty"`
    DateEnd     *time.Time `json:"date_end,omitempty"`
    CheckPoint  int        `json:"check_point"`
    ImgPath     string     `json:"img_path,omitempty"`
}

type ThemeUpdateRequest struct {
    Name        string     `json:"name" validate:"required,min=3"`
    Description string     `json:"description,omitempty"`
    IsActive    bool       `json:"is_active"`
    DateBegin   *time.Time `json:"date_begin,omitempty"`
    DateEnd     *time.Time `json:"date_end,omitempty"`
    CheckPoint  int        `json:"check_point"`
    ImgPath     string     `json:"img_path,omitempty"`
}

type ThemeResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    IsActive    bool   `json:"is_active"`
    CreatedByID string `json:"created_by_id"`
    CreatedAt   string `json:"created_at"`
    DateBegin   string `json:"date_begin"`
    DateEnd     string `json:"date_end"`
    MaxPoint    int    `json:"max_point"`
    CheckPoint  int    `json:"check_point"`
    ImgPath     string `json:"img_path"`
}

type ThemeWithDetailResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    IsActive    bool   `json:"is_active"`
    CreatedByID string `json:"created_by_id"`
    CreatedAt   string `json:"created_at"`
    DateBegin   string `json:"date_begin"`
    DateEnd     string `json:"date_end"`
    MaxPoint    int    `json:"max_point"`
    CheckPoint  int    `json:"check_point"`
    ImgPath     string `json:"img_path"`
    CreatedBy   UserShortResponse `json:"created_by"`
}

func ThemeToResponse(theme *domain.ThemeWithDetail) *ThemeWithDetailResponse {
    var dateBegin, dateEnd string
    if theme.DateBegin != nil {
        dateBegin = theme.DateBegin.Format(time.RFC3339)
    }
    if theme.DateEnd != nil {
        dateEnd = theme.DateEnd.Format(time.RFC3339)
    }

    return &ThemeWithDetailResponse{
        ID:          theme.ID.String(),
        Name:        theme.Name,
        Description: theme.Description,
        IsActive:    theme.IsActive,
        CreatedByID: theme.CreatedByID.String(),
        CreatedAt:   theme.CreatedAt.Format(time.RFC3339),
        DateBegin:   dateBegin,
        DateEnd:     dateEnd,
        MaxPoint:    theme.MaxPoint,
        CheckPoint:  theme.CheckPoint,
        ImgPath:     theme.ImgPath,
        CreatedBy: UserShortResponse{
            ID:   theme.CreatedBy.ID.String(),
            Name: theme.CreatedBy.Name,
        },
    }
}

func ThemesToResponseSlice(themes []*domain.ThemeWithDetail) []*ThemeWithDetailResponse {
    if themes == nil {
        return nil
    }
    out := make([]*ThemeWithDetailResponse, len(themes))
    for i, item := range themes {
        out[i] = ThemeToResponse(item)
    }
    return out
}