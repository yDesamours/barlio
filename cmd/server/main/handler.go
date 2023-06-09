package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/helper"
	"barlio/internal/validator"
	"barlio/ui"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

func (app *App) homePageHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserHelper(r)
	data := app.newTemplateData(user, r)

	app.setLastPageHelper(r)

	tmpl := app.PageTemplates.Get(HOMEPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(nil, r)
	data.Set("showHeader", false)

	tmpl := app.PageTemplates.Get(SIGNINPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signupHandler(w http.ResponseWriter, r *http.Request) {
	form := app.readFormDataHelper(r)
	infos := app.newTemplateData(nil, r)

	user := app.newUser(form)

	validator := validator.New()
	app.validateUser(user, form.Get("passwordconfirm"), validator)
	if !validator.Valid() {
		err := app.signInError(w, infos, form, validator)
		if err != nil {
			app.error(err)
		}
		return
	}

	err := app.validateUserUnicity(user, validator)
	if err != nil {
		app.error(err)
		return
	}

	hash, err := helper.HashPassword(user.Password)
	if err != nil {
		app.error(err)
		return
	}
	user.Password = hash

	err = app.models.user.Insert(user)
	if err != nil {
		app.error(err)
		return
	}

	app.SessionManager.Put(r.Context(), "userId", user.ID)

	http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
}

func (app *App) emailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(nil, r)
	data.Set("showHeader", false)

	user := app.getUserHelper(r)
	if user == nil {
		http.Redirect(w, r, SIGNUPPAGE, http.StatusMovedPermanently)
		return
	}

	app.SessionManager.Put(r.Context(), "lastpage", EMAILVERIFICATIONPAGE)

	tmpl := app.PageTemplates.Get(EMAILVERIFICATIONPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.token.DeleteForUser(user.ID, model.EMAILVERIFICATIONSCOPE)
	if err != nil {
		app.error(err)
		return
	}

	token, err := app.newVerificationToken(user)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.token.Insert(token)
	if err != nil {
		app.error(err)
		return
	}

	go app.sendVerificationEmail(user, token)
}

func (app *App) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(nil, r)
	app.setLastPageHelper(r)
	data.Set("showHeader", false)

	app.SessionManager.Put(r.Context(), "lastpage", SIGNUPPAGE)

	tmpl := app.PageTemplates.Get(SIGNUPPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) signinHandler(w http.ResponseWriter, r *http.Request) {
	var lastPage string
	loginUser := app.newUser(app.readFormDataHelper(r))
	templateData := templateData{"username": loginUser.Username, "password": ""}

	user, err := app.models.user.Get(*loginUser)
	if errors.Is(err, sql.ErrNoRows) {
		templateData["errors"] = map[string][]string{"username": {"username not found"}}
		tmpl := app.PageTemplates.Get(SIGNINPAGE)
		err := tmpl.Execute(w, templateData)
		if err != nil {
			app.error(err)
		}
		return
	}

	if same := helper.CompareHash(loginUser.Password, user.Password); !same {
		templateData["errors"] = map[string][]string{"password": {"incorrect password"}}
		tmpl := app.PageTemplates.Get(SIGNINPAGE)
		err := tmpl.Execute(w, templateData)
		if err != nil {
			app.error(err)
		}
		return
	}

	app.SessionManager.Put(r.Context(), "userId", user.ID)
	if !user.IsVerified {
		http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
		return
	}

	if lastPage = app.SessionManager.GetString(r.Context(), "lastpage"); lastPage == "" {
		lastPage = HOMEPAGE
	}
	http.Redirect(w, r, lastPage, http.StatusMovedPermanently)
}

func (app *App) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")

	if tokenString == "" {
		http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
		return
	}

	user := app.getUserHelper(r)
	if user == nil {
		http.Redirect(w, r, SIGNUPPAGE, http.StatusMovedPermanently)
		return
	}

	userToken, err := app.models.token.GetForUser(user.ID, model.EMAILVERIFICATIONSCOPE)
	if errors.Is(err, sql.ErrNoRows) {
		http.Redirect(w, r, EMAILVERIFICATIONPAGE, http.StatusMovedPermanently)
		return
	}

	validator := validator.New()
	app.validateTokenHelper(userToken, tokenString, validator)
	if !validator.Valid() {
		validator.Error()
		return
	}

	err = app.models.user.Activate(user)
	if err != nil {
		app.error(err)
		return
	}

	app.SessionManager.Put(r.Context(), "flash", "Account Verified")
	http.Redirect(w, r, HOMEPAGE, http.StatusMovedPermanently)
}

func (app *App) logoutHandler(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(nil, r)
	data.Set("page", HOMEPAGE)

	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.error(err)
		return
	}
	app.SessionManager.Clear(r.Context())
	http.Redirect(w, r, HOMEPAGE, http.StatusMovedPermanently)
}

func (app *App) updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	form := app.readFormDataHelper(r)
	user := app.getUserHelper(r)
	tmpl := app.PageTemplates.Get(PROFILEPAGE)

	app.updateUserProfile(user, form)

	err := app.models.user.UpdateUser(user)
	if err != nil {
		app.error(err)
	}

	data := app.newTemplateData(user, r)
	err = tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) securityPageHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserHelper(r)
	data := app.newTemplateData(user, r)
	tmpl := app.PageTemplates.Get(SECURITYPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) changeUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	form := app.readFormDataHelper(r)
	user := app.getUserHelper(r)
	validator := validator.New()
	data := app.newTemplateData(user, r)
	tmpl := app.PageTemplates.Get(SECURITYPAGE)

	app.validateChangeUserPasswordForm(form, validator)
	if !validator.Valid() {
		data["errors"] = validator.Error()
		tmpl.Execute(w, data)
		return
	}

	token, err := app.newPasswordChangeToken(user)
	if err != nil {
		app.error(err)
		return
	}

	request, err := app.newPasswordChangeRequest(user, form, token)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.token.Insert(token)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.request.Insert(request)
	if err != nil {
		app.error(err)
		return
	}

	data.Set("token", token.Token)
	tmpl.Execute(w, data)

}

func (app *App) confirmPasswordChangeHandler(w http.ResponseWriter, r *http.Request) {
	form := app.readFormDataHelper(r)
	user := app.getUserHelper(r)
	validator := validator.New()

	changePasswordStruct := app.changePasswordConfirmStruct(form)

	token, err := app.models.token.GetForUser(user.ID, model.PASSWORDCHANGESCOPE)
	if err != nil {
		app.error(err)
		return
	}

	app.validateTokenHelper(token, string(changePasswordStruct.token), validator)
	if !validator.Valid() {
		return
	}

	request, err := app.models.request.GetForUser(user.ID, model.PASSWORDCHANGESCOPE)
	if err != nil {
		app.error(err)
		return
	}

	var cpRequest changePasswordRequest
	err = json.Unmarshal(request.Data, &cpRequest)
	if err != nil {
		app.error(err)
		return
	}

	app.models.request.DeleteForUser(user.ID, model.PASSWORDCHANGESCOPE)
	if err != nil {
		app.error(err)
		return
	}

	app.models.token.DeleteForUser(user.ID, model.PASSWORDCHANGESCOPE)
	if err != nil {
		app.error(err)
		return
	}

	err = app.models.user.UpdateUserPassword(user)
	if err != nil {
		app.error(err)
		return
	}

	http.Redirect(w, r, SECURITYPAGE, http.StatusMovedPermanently)

}

func (app *App) notFound(w http.ResponseWriter, r *http.Request) {
	notFoundPath := filepath.Join("web", "html", "notFound.tmpl.html")
	temp, err := template.ParseFiles(notFoundPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	temp.Execute(w, templateData{})
}

func (app *App) fileServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = fmt.Sprintf("/web%s", r.URL.Path)
		server := http.FileServer(http.FS(ui.FILES))
		server.ServeHTTP(w, r)
	})
}

func (app *App) docServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastindex := strings.LastIndex(r.URL.Path, "/")
		r.URL.Path = r.URL.Path[lastindex+1:]
		server := http.FileServer(http.Dir("C:\\Users\\User\\code\\go\\src\\barlio\\docs"))
		server.ServeHTTP(w, r)
	})
}

func (app *App) profilePageHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserHelper(r)
	data := app.newTemplateData(user, r)

	links := []string{
		"https://cdnjs.cloudflare.com/ajax/libs/cropperjs/1.5.12/cropper.min.css",
		"https://cdn.jsdelivr.net/npm/cropperjs@1.5.11/dist/cropper.min.css",
	}
	scripts := []string{
		"/statics/js/profilphoto.js",
		"https://cdnjs.cloudflare.com/ajax/libs/cropperjs/1.5.12/cropper.min.js",
	}
	data.Set("scripts", scripts)
	data.Set("links", links)
	data.Set("profilPhoto", "/docs/defaultavatar.jpg")

	app.SessionManager.Put(r.Context(), "lastpage", PROFILEPAGE)

	tmpl := app.PageTemplates.Get(PROFILEPAGE)
	err := tmpl.Execute(w, data)
	if err != nil {
		app.error(err)
	}
}

func (app *App) updateUserProfilePicture(w http.ResponseWriter, r *http.Request) {
	user := app.getUserHelper(r)
	form, err := app.readMultipartFormDataHelper(r)
	if err != nil {
		app.error(err)
		return
	}

	file, err := app.readMultipartFile(form, "file")
	if err != nil {
		app.error(err)
		return
	}

	err = app.saveFile(string(user.ProfilPicture), file)
	if err != nil {
		app.error(err)
		return
	}
}
