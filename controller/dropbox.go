package controller

import (
	"fmt"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
	"lens-locked-go/config"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"net/http"
	"time"
)

type dropboxController struct {
	dropboxCfg     *oauth2.Config
	dropboxService service.IDropboxService
}

func NewDropboxController(cfg *config.Config, ds service.IDropboxService) *dropboxController {
	dropboxCfg := &oauth2.Config{
		ClientID:     cfg.Dropbox.Id,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.Dropbox.AuthUrl,
			TokenURL: cfg.Dropbox.TokenUrl,
		},
		RedirectURL: cfg.Dropbox.RedirectUrl,
	}

	return &dropboxController{
		dropboxCfg:     dropboxCfg,
		dropboxService: ds,
	}
}

func (c *dropboxController) ConnectGet(w http.ResponseWriter, req *http.Request) {
	state := csrf.Token(req)

	cookie := &http.Cookie{
		Name:     "dropbox_token",
		Value:    state,
		Expires:  time.Now().Add(time.Minute * 5),
		HttpOnly: true,
	}

	url := c.dropboxCfg.AuthCodeURL(state)

	http.SetCookie(w, cookie)

	http.Redirect(w, req, url, http.StatusFound)
}

func (c *dropboxController) CallbackGet(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		//Change status code
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	state := req.FormValue("state")

	cookie, err := req.Cookie("dropbox_token")

	if err != nil {
		//Change status code
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if state != cookie.Value {
		//Change status code
		http.Error(w, "invalid state", http.StatusInternalServerError)

		return
	}
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)

	http.SetCookie(w, cookie)

	code := req.FormValue("code")

	token, err := c.dropboxCfg.Exchange(req.Context(), code)

	if err != nil {
		//Change status code
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	user, er := context.User(req.Context())

	if er != nil {
		http.Error(w, er.Message, er.StatusCode)

		return
	}
	existingDropbox, er := c.dropboxService.GetByUserId(user.ID)

	if er != nil && er.StatusCode != http.StatusNotFound {
		http.Error(w, er.Message, er.StatusCode)

		return
	}

	if existingDropbox != nil {
		er = c.dropboxService.Delete(existingDropbox)

		if er != nil {
			http.Error(w, er.Message, er.StatusCode)

			return
		}
	}
	dropbox := &model.Dropbox{
		Base:   model.Base{},
		UserId: user.ID,
		Token:  *token,
	}

	er = c.dropboxService.Create(dropbox)

	if er != nil {
		http.Error(w, er.Message, er.StatusCode)

		return
	}
	fmt.Printf("%+v", token)
}

func (c *dropboxController) ConnectRoute() string {
	return dropboxConnectRoute
}

func (c *dropboxController) CallbackRoute() string {
	return dropboxCallbackRoute
}
