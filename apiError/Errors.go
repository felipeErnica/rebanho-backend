package apiError

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    int    `json:"-"`
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

const InternalError = "InternalError"
const ConflictError = "ConflictError"
const IncorretInfoError = "IncorretInfoError"
const OtherError = "OtherError"
const ApiWarning = "ApiWarning"

func NewAPIError(code int, kind string, title string, message string) *APIError {
	return &APIError{
		Code:    code,
		Kind:    kind,
		Title:   "ERRO: " + title,
		Message: message,
	}
}

func ConflictAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		Kind:    ConflictError,
		Title:   "ERRO: Informação já existe!",
		Message: message,
	}
}

func IncorrectEntityAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUnprocessableEntity,
		Kind:    IncorretInfoError,
		Title:   "ERRO: Informação Incorreta!",
		Message: message,
	}
}

func InternalServerAPIError(err error) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Kind:    InternalError,
		Title:   "GRAVE: Erro Interno do Servidor!",
		Message: fmt.Sprintf("Erro causado por: %s. Consulte o suporte para resolver esta questão.", err.Error()),
	}
}

func ConflictAPIWarning(message string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		Kind:    ApiWarning,
		Title:   "AVISO: Informação já existe!",
		Message: message,
	}
}

func NewAPIWarning(message string, title string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		Kind:    ApiWarning,
		Title:   "AVISO: ",
		Message: message,
	}
}

func EmptyList() error {
	err := errors.New("A matriz está vazia")
	return err
}
