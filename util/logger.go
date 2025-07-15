package util

import (
	"fmt"
	"log/slog"
	"os"
)

var handler = slog.NewTextHandler(os.Stdout, nil)
var logger = slog.New(handler)

func LogDomainsInit(name string) {
    LogInfo("Domínio de " + name + " iniciado com sucesso!", false)
}

func LogInfo(msg string, jumpLine bool) {
    if jumpLine {
        fmt.Println()
    }
    logger.Info(msg)
}

func LogDebug(msg string) {
    logger.Info(msg)
}

func LogError(msg string) {
    logger.Error(msg)
}
