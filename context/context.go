package context

import (
	"context"
	"lens-locked-go/model"
)

type contextKey string

const key contextKey = "user"

const noUserInContextErrorMessage = "no user in context"

func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, key, user)
}

func User(ctx context.Context) (*model.User, *model.Error) {
	value := ctx.Value(key)

	if value != nil {
		if user, ok := value.(*model.User); ok {
			return user, nil
		}
	}
	return nil, model.NewInternalServerApiError(noUserInContextErrorMessage)
}
