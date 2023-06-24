package main

import "barlio/internal/types"

const (
	userType types.String = "user"
)

type changePasswordRequest struct {
	Password string
}

type changePasswordConfirm struct {
	Password types.String
	token    types.String
}
