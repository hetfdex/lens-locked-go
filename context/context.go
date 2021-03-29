package context

import (
	"context"
	"lens-locked-go/model"
)

const userKey = "user"
const noUserInContextErrorMessage = "no user in context"

func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) (*model.User, *model.Error) {
	value := ctx.Value(userKey)

	if value != nil {
		if user, ok := value.(*model.User); ok {
			return user, nil
		}
	}
	return nil, model.NewInternalServerApiError(noUserInContextErrorMessage)
}
