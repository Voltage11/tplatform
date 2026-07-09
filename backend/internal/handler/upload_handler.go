package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Voltage11/tplatform/internal/config"
	"github.com/Voltage11/tplatform/internal/handler/httputils"
	"github.com/Voltage11/tplatform/internal/middleware"
	"github.com/Voltage11/tplatform/internal/types/apperror"
	"github.com/Voltage11/tplatform/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UploadHandler struct {
    cfg    config.UploadConfig
    logger logger.Logger
}

func NewUploadHandler(r chi.Router, authMW *middleware.AuthMiddleware, cfg config.UploadConfig, log logger.Logger) {
    h := &UploadHandler{cfg: cfg, logger: log}

    // Защищаем маршрут – только авторизованные пользователи
    r.Group(func(r chi.Router) {
        r.Use(authMW.RequireAuth)
        r.Post("/api/v1/upload", h.Upload)
    })
}

func (h *UploadHandler) Upload(w http.ResponseWriter, r *http.Request) {
    // Ограничиваем размер тела запроса
    r.Body = http.MaxBytesReader(w, r.Body, h.cfg.MaxSizeMB<<20) 

    // Парсим форму
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        httputils.WriteError(w, apperror.NewBadRequest("файл слишком большой или неверный формат", err))
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        httputils.WriteError(w, apperror.NewBadRequest("не удалось получить файл", err))
        return
    }
    defer file.Close()

    // Проверяем тип файла
    contentType := header.Header.Get("Content-Type")
    allowedTypes := map[string]bool{
        "image/jpeg":    true,
        "image/png":     true,
        "image/gif":     true,
        "image/svg+xml": true,
    }
    if !allowedTypes[contentType] {
        httputils.WriteError(w, apperror.NewBadRequest("недопустимый тип файла. Разрешены JPEG, PNG, GIF, SVG", nil))
        return
    }

    // Проверяем расширение
    ext := strings.ToLower(filepath.Ext(header.Filename))
    if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".svg" {
        httputils.WriteError(w, apperror.NewBadRequest("недопустимое расширение файла", nil))
        return
    }

    // Создаём директорию, если её нет
    if err := os.MkdirAll(h.cfg.Dir, 0755); err != nil {
        httputils.WriteError(w, apperror.NewInternal("ошибка создания директории", err))
        return
    }

    // Генерируем уникальное имя файла
    filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
    destPath := filepath.Join(h.cfg.Dir, filename)

    // Сохраняем файл
    dst, err := os.Create(destPath)
    if err != nil {
        httputils.WriteError(w, apperror.NewInternal("не удалось создать файл", err))
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        httputils.WriteError(w, apperror.NewInternal("ошибка сохранения файла", err))
        return
    }

    // Формируем URL для доступа (относительный путь)
    fileURL := fmt.Sprintf("/uploads/%s", filename)

    h.logger.Info("Файл загружен", "file", fileURL)

    // Возвращаем URL
    httputils.WriteJSON(w, http.StatusCreated, map[string]string{"url": fileURL})
}