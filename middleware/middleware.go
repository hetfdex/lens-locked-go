package middleware

import (
	"lens-locked-go/context"
	"lens-locked-go/controller"
	"lens-locked-go/service"
	"net/http"
)

type Middleware struct {
	service.IUserService
}

func NewMiddleware(us service.IUserService) *Middleware {
	return &Middleware{
		us,
	}
}

func (m *Middleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(controller.CookieName)

		if err != nil {
			controller.Redirect(w, req, controller.LoginUserRoute)

			return
		}
		user, er := m.LoginWithToken(cookie.Value)

		if er != nil {
			controller.Redirect(w, req, controller.LoginUserRoute)

			return
		}
		ctx := req.Context()

		ctx = context.WithUser(ctx, user)

		req = req.WithContext(ctx)

		next(w, req)
	}
}
