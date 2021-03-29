package middleware

import (
	"lens-locked-go/controller"
	"lens-locked-go/service"
	"log"
	"net/http"
)

type RequireUser struct {
	service.IUserService
}

func NewRequireUserMiddleware(us service.IUserService) *RequireUser {
	return &RequireUser{
		us,
	}
}

func (m *RequireUser) Apply(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(controller.CookieName)

		if err != nil {
			controller.Redirect(w, req, "/login")

			return
		}
		user, er := m.LoginWithToken(cookie.Value)

		if er != nil {
			controller.Redirect(w, req, "/login")

			return
		}
		log.Println(user.ID, user.Email)

		next(w, req)
	})
}
