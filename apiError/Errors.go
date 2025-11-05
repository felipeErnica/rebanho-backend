package apiError

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    int    `json:"-"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewAPIError(code int, title string, message string) *APIError {
	return &APIError{
		Code:    code,
		Title:   "ERRO: " + title,
		Message: message,
	}
}

func ConflictAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		Title:   "ERRO: Informação já existe!",
		Message: message,
	}
}

func IncorrectEntityAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUnprocessableEntity,
		Title:   "ERRO: Informação Incorreta!",
		Message: message,
	}
}

func InternalServerAPIError(err error) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Title: "GRAVE: Erro Interno do Servidor!",
		Message: fmt.Sprintf("Erro causado por: %s. Consulte o suporte para resolver esta questão.", err.Error()),
	}
}

func EmptyList() error {
	err := errors.New("A matriz está vazia")
	return err
}
