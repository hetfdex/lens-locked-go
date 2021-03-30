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

func (m *Middleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, _ := context.User(req.Context())

		if user != nil {
			next.ServeHTTP(w, req)

			return
		}
		cookie, err := req.Cookie(controller.CookieName())

		if err != nil {
			next.ServeHTTP(w, req)

			return
		}
		user, er := m.LoginWithToken(cookie.Value)

		if er != nil {
			next.ServeHTTP(w, req)

			return
		}
		ctx := req.Context()

		ctx = context.WithUser(ctx, user)

		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func (m *Middleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user, _ := context.User(req.Context())

		if user == nil {
			controller.Redirect(w, req, controller.LoginUserRoute())

			return
		}
		next(w, req)
	}
}
