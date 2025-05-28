package util

import (
	"log/slog"
	"os"
)

var handler = slog.NewTextHandler(os.Stdout, nil)
var logger = slog.New(handler)

func LogDomainsInit(name string) {
    LogInfo("Domínio de " + name + " iniciado com sucesso!")
}

func LogInfo(msg string) {
    logger.Info(msg)
}

func LogDebug(msg string) {
    logger.Info(msg)
}

func LogError(msg string) {
    logger.Error(msg)
}
