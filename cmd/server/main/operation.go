package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/token"
	"barlio/internal/types"
	"barlio/internal/validator"
	"bytes"
	"net/http"
	"net/url"
	"time"
)

const (
	emailVerificationTokenDuration = time.Hour * 4
	passwordChangeTokenDuration    = time.Minute
)

func (app *App) signInError(w http.ResponseWriter, data templateData, form url.Values, validator *validator.Validator) error {
	signinErrors := validator.Error()
	data.Set("errors", signinErrors)
	data.Set("signin", map[string]interface{}{"username": form.Get("username"), "email": form.Get("email")})
	data.Set("page", "Signin")

	tmpl := app.PageTemplates["signin"]
	return tmpl.Execute(w, data)
}

func (app *App) sendVerificationEmail(user *model.User, token *model.Token) {
	var (
		data         = templateData{"token": token.Token}
		mailTemplate = app.MailTemplate.Get("/emailverification")
	)

	mailObject, err := parseEmailData(mailTemplate, data)
	if err != nil {
		app.error(err)
		return
	}
	go app.Mailer.Send(string(user.Email), mailObject)
}

func parseEmailData(mailTemplate *PageTemplate, data templateData) (map[string]string, error) {
	var (
		buffer     = &bytes.Buffer{}
		mailObject = map[string]string{}
	)

	err := mailTemplate.ExecuteTemplate(buffer, "subject", data)
	if err != nil {
		return nil, err
	}
	mailObject["subject"] = buffer.String()
	buffer.Reset()

	err = mailTemplate.ExecuteTemplate(buffer, "alternative", data)
	if err != nil {
		return nil, err
	}
	mailObject["alternative"] = buffer.String()
	buffer.Reset()

	err = mailTemplate.ExecuteTemplate(buffer, "body", data)
	if err != nil {
		return nil, err
	}
	mailObject["body"] = buffer.String()

	return mailObject, nil
}

func (app *App) newVerificationToken(user *model.User) (*model.Token, error) {
	plainTextToken, hashedToken, err := token.GenerateToken()
	if err != nil {
		return nil, err
	}
	token := model.Token{
		UserId:    user.ID,
		Scope:     model.EMAILVERIFICATIONSCOPE,
		Token:     plainTextToken,
		Hash:      hashedToken,
		ExpiretAt: time.Now().Add(emailVerificationTokenDuration),
	}
	return &token, nil
}

func (app *App) updateUserProfilInfos(user *model.User, form url.Values) {
	user.Firstname = types.String(form.Get("firstname"))
	user.Lastname = types.String(form.Get("lastname"))
	user.Bio = types.String(form.Get("bio"))
	birthdate, _ := time.Parse(time.DateOnly, form.Get("birhtdate"))
	user.Birthdate = birthdate
}

func (app *App) newPasswordChangeToken(user *model.User) (*model.Token, error) {
	plainTextToken, hashedToken, err := token.GenerateToken()
	if err != nil {
		return nil, err
	}

	token := model.Token{
		UserId:    user.ID,
		Scope:     model.PASSWORDCHANGESCOPE,
		Token:     plainTextToken,
		Hash:      hashedToken,
		ExpiretAt: time.Now().Add(passwordChangeTokenDuration),
	}
	return &token, nil
}

func (app *App) changePasswordConfirmStruct(form url.Values) *changePasswordConfirm {
	return &changePasswordConfirm{Password: types.String(form.Get("password")), token: types.String(form.Get("token"))}
}
