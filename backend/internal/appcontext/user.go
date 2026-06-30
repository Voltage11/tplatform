package appcontext

import (
	"context"

	"github.com/Voltage11/tplatform/internal/domain"
)

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

func GetUserFromContext(ctx context.Context) *domain.User {
	user, ok := ctx.Value(UserContextKey).(*domain.User)
	if !ok {
		return nil
	}
	return user
}

func SetUserToContext(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}