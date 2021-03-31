package middleware

import (
	"lens-locked-go/context"
	"lens-locked-go/service"
	"lens-locked-go/util"
	"net/http"
)

type Middleware struct {
	userService service.IUserService
	loginRoute  string
}

func NewMiddleware(us service.IUserService, lr string) *Middleware {
	return &Middleware{
		userService: us,
		loginRoute:  lr,
	}
}

func (m *Middleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, _ := context.User(req.Context())

		if user != nil {
			next.ServeHTTP(w, req)

			return
		}
		cookie, err := req.Cookie(util.CookieName)

		if err != nil {
			next.ServeHTTP(w, req)

			return
		}
		user, er := m.userService.LoginWithToken(cookie.Value)

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
			util.Redirect(w, req, m.loginRoute)

			return
		}
		next(w, req)
	}
}
