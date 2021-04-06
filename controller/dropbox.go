package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
	"io"
	"lens-locked-go/config"
	"lens-locked-go/context"
	"lens-locked-go/model"
	"lens-locked-go/service"
	"log"
	"net/http"
	"time"
)

const dropboxCookieName = "dropbox_token"

const invalidAuthenticationErrorMessage = "invalid authentication"

type dropboxController struct {
	cfg            *oauth2.Config
	dropboxService service.IDropboxService
}

func NewDropboxController(dropboxCfg *config.DropboxConfig, ds service.IDropboxService) *dropboxController {
	cfg := &oauth2.Config{
		ClientID:     dropboxCfg.Id,
		ClientSecret: dropboxCfg.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  dropboxCfg.AuthUrl,
			TokenURL: dropboxCfg.TokenUrl,
		},
		RedirectURL: dropboxCfg.RedirectUrl,
	}

	return &dropboxController{
		cfg:            cfg,
		dropboxService: ds,
	}
}

func (c *dropboxController) ConnectGet(w http.ResponseWriter, req *http.Request) {
	state := csrf.Token(req)

	cookie := &http.Cookie{
		Name:     dropboxCookieName,
		Value:    state,
		Expires:  time.Now().Add(time.Minute * 5),
		HttpOnly: true,
	}

	url := c.cfg.AuthCodeURL(state)

	http.SetCookie(w, cookie)

	http.Redirect(w, req, url, http.StatusFound)
}

func (c *dropboxController) CallbackGet(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	state := req.FormValue("state")

	cookie, err := req.Cookie(dropboxCookieName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	if state != cookie.Value {
		http.Error(w, invalidAuthenticationErrorMessage, http.StatusUnauthorized)

		return
	}
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)

	http.SetCookie(w, cookie)

	code := req.FormValue("code")

	token, err := c.cfg.Exchange(req.Context(), code)

	if err != nil {
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
	log.Printf("%+v", token)
}

func (c *dropboxController) QueryGet(w http.ResponseWriter, req *http.Request) {
	er := req.ParseForm()

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)

		return
	}
	path := req.FormValue("path")

	user, err := context.User(req.Context())

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	dropbox, err := c.dropboxService.GetByUserId(user.ID)

	if err != nil {
		http.Error(w, err.Message, err.StatusCode)

		return
	}
	token := dropbox.Token

	client := c.cfg.Client(req.Context(), &token)

	data := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}

	dataBytes, er := json.Marshal(&data)

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)

		return
	}
	r, er := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder", bytes.NewReader(dataBytes))

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)

		return
	}
	r.Header.Add("Content-Type", "application/json")

	res, er := client.Do(r)

	if er != nil {
		http.Error(w, er.Error(), http.StatusFailedDependency)

		return
	}
	defer res.Body.Close()

	_, er = io.Copy(w, res.Body)

	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)

		return
	}
	log.Println(res.StatusCode)
}

func (c *dropboxController) ConnectRoute() string {
	return dropboxConnectRoute
}

func (c *dropboxController) CallbackRoute() string {
	return dropboxCallbackRoute
}

func (c *dropboxController) QueryRoute() string {
	return queryRoute
}
