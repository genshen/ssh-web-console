package models

import "github.com/genshen/ssh-web-console/src/utils"

const (
	SIGN_IN_FORM_TYPE_ERROR_VALID    = iota
	SIGN_IN_FORM_TYPE_ERROR_PASSWORD
	SIGN_IN_FORM_TYPE_ERROR_TEST
)

type UserInfo struct {
	utils.Connection
	Username string `json:"username"`
	Password string `json:"-"`
}

type JsonResponse struct {
	HasError bool        `json:"has_error"`
	Message  interface{} `json:"message"`
	Addition interface{} `json:"addition"`
}
