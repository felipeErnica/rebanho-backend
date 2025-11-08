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

const INTERNAL_ERROR = "InternalError"
const CONFLICT_ERROR = "ConflictError"
const INFO_INCORRET_ERROR = "IncorretInfoError"
const DELETE_ERROR = "DeleteError"
const OTHER_ERROR = "OtherError"
const API_WARNING = "ApiWarning"

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
		Kind:    CONFLICT_ERROR,
		Title:   "ERRO: Informação já existe!",
		Message: message,
	}
}

func IncorrectEntityAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUnprocessableEntity,
		Kind:    INFO_INCORRET_ERROR,
		Title:   "ERRO: Informação Incorreta!",
		Message: message,
	}
}

func InternalServerAPIError(err error) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Kind:    INTERNAL_ERROR,
		Title:   "GRAVE: Erro Interno do Servidor!",
		Message: fmt.Sprintf("Erro causado por: %s. Consulte o suporte para resolver esta questão.", err.Error()),
	}
}

func DeleteAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUpgradeRequired,
		Kind:    DELETE_ERROR,
		Title:   "AVISO: Este objeto possui dependências.",
		Message: message,
	}
}

func ConflictAPIWarning(message string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		Kind:    API_WARNING,
		Title:   "AVISO: Informação já existe!",
		Message: message,
	}
}

func DeleteWarning(message string) *APIError {
	return &APIError{
		Code:    http.StatusUpgradeRequired,
		Kind:    API_WARNING,
		Title:   "AVISO: Este objeto possui dependências.",
		Message: message,
	}
}

func NewAPIWarning(message string, title string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		Kind:    API_WARNING,
		Title:   "AVISO: ",
		Message: message,
	}
}

func EmptyList() error {
	err := errors.New("A matriz está vazia")
	return err
}
