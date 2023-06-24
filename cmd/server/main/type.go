package main

import "barlio/internal/types"

const (
	userType types.String = "user"
)

type changePassword struct {
	Password       types.String
	PassordConfirm types.String
}
