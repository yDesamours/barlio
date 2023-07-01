package main

import (
	"barlio/cmd/server/model"
	"barlio/internal/helper"
	"barlio/internal/token"
	"barlio/internal/types"
	"barlio/internal/validator"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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

func (app *App) updateUserProfile(user *model.User, form url.Values) {
	user.Firstname = types.String(form.Get("firstname"))
	user.Lastname = types.String(form.Get("lastname"))
	user.Bio = types.String(form.Get("bio"))
	birthdate, _ := time.Parse(time.DateOnly, form.Get("birthdate"))

	user.Birthdate = birthdate
}

func (app *App) saveFile(filename string, content []byte) error {
	return os.WriteFile(filename, content, os.ModePerm)
}

func (app *App) readMultipartFile(form *multipart.Form, filename string) ([]byte, error) {
	headers := form.File[filename]
	if len(headers) == 0 {
		return nil, fmt.Errorf("no file submitted")
	}

	header := headers[0]
	file, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content := make([]byte, header.Size)
	_, err = file.Read(content)
	return content, nil
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

func (app *App) newUser(form url.Values) *model.User {
	user := &model.User{
		Username: types.String(form.Get("username")),
		Email:    types.String(form.Get("email")),
		Password: types.String(form.Get("password")),
	}
	return user
}

func (app *App) newPasswordChangeRequest(user *model.User, form url.Values, token *model.Token) (*model.Request, error) {
	password, err := helper.HashPassword(form.Get("password"))
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(map[string]string{"password": password})
	if err != nil {
		return nil, err
	}

	return &model.Request{
		UserId:    user.ID,
		Scope:     model.PASSWORDCHANGESCOPE,
		CreatedAt: time.Now(),
		ExpireAt:  token.ExpiretAt,
		Data:      b,
	}, nil
}

func (app *App) changePasswordConfirmStruct(form url.Values) *changePasswordConfirm {
	return &changePasswordConfirm{Password: types.String(form.Get("password")), token: types.String(form.Get("token"))}
}
