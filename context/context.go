package context

import (
	"context"
	"lens-locked-go/model"
)

const userKey = "user"

func WithUser(ctx context.Context, user *model.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *model.User {
	value := ctx.Value(userKey)

	if value != nil {
		if user, ok := value.(*model.User); ok {
			return user
		}
	}
	return nil
}
