package apiError

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    int    `json:"-"`
	Kind    string `json:"kind"`
	ErrType string `json:"errType"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

const ERROR_TYPE = "Error"
const WARNING_TYPE = "Warning"

const INTERNAL_ERROR = "InternalError"
const CONFLICT_ERROR = "ConflictError"
const INFO_INCORRET_ERROR = "IncorretInfoError"
const DELETE_ERROR = "DeleteError"

const CONFLICT_WARNING = "ConflictWarning"
const TRANSFER_WARNING = "TransferWarning"

func NewAPIError(code int, kind string, title string, message string) *APIError {
	return &APIError{
		Code:    code,
		ErrType: ERROR_TYPE,
		Kind:    kind,
		Title:   "ERRO: " + title,
		Message: message,
	}
}

func ConflictAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		ErrType: ERROR_TYPE,
		Kind:    CONFLICT_ERROR,
		Title:   "ERRO: Informação já existe!",
		Message: message,
	}
}

func IncorrectEntityAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUnprocessableEntity,
		ErrType: ERROR_TYPE,
		Kind:    INFO_INCORRET_ERROR,
		Title:   "ERRO: Informação Incorreta!",
		Message: message,
	}
}

func InternalServerAPIError(err error) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		ErrType: ERROR_TYPE,
		Kind:    INTERNAL_ERROR,
		Title:   "GRAVE: Erro Interno do Banco de Dados!",
		Message: fmt.Sprintf("Erro causado por: %s. Consulte o suporte para resolver esta questão.", err.Error()),
	}
}

func DeleteAPIError(message string) *APIError {
	return &APIError{
		Code:    http.StatusUpgradeRequired,
		ErrType: ERROR_TYPE,
		Kind:    DELETE_ERROR,
		Title:   "AVISO: Este objeto possui dependências.",
		Message: message,
	}
}

func ConflictAPIWarning(message string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		ErrType: WARNING_TYPE,
		Kind:    CONFLICT_WARNING,
		Title:   "AVISO: Informação já existe!",
		Message: message,
	}
}

func DeleteWarning(message string) *APIError {
	return &APIError{
		Code:    http.StatusUpgradeRequired,
		ErrType: WARNING_TYPE,
		Kind:    CONFLICT_WARNING,
		Title:   "AVISO: Este objeto possui dependências.",
		Message: message,
	}
}

func DeleteWarningKind(kind string, message string) *APIError {
	return &APIError{
		Code:    http.StatusUpgradeRequired,
		ErrType: WARNING_TYPE,
		Kind:    kind,
		Title:   "AVISO: Este objeto possui dependências.",
		Message: message,
	}
}

func NewAPIWarning(title string, message string, kind string) *APIError {
	return &APIError{
		Code:    http.StatusConflict,
		ErrType: WARNING_TYPE,
		Kind:    kind,
		Title:   "AVISO: " + title,
		Message: message,
	}
}

func EmptyList() error {
	err := errors.New("A matriz está vazia")
	return err
}
