package util

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func writeResponse(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

/*Retorna uma matriz baseado no texto enviado, usa a vírgula por separador comum.
Se o texto estiver vazio, retorna uma matriz vazia.*/
func ParseArray(inputString string) []string {
    if inputString == "" {
        return []string{}
    }
    array := strings.Split(inputString, ",")
    return array
}

func ParseBool(str string) (bool, error) {
	if str == "" {
		return false, nil
	}
	bln, err := strconv.ParseBool(str)
	return bln, err
}

func FormatWarningMessage(messages ...string) string {
	msgBody := FormatMessageBody(messages...)
	return "As seguintes ocorrências foram detectadas:" + msgBody
}

func FormatMessageBody(messages ...string) string {
	formatedMsg := []string{}
	for i, message := range messages {
		msg := fmt.Sprintf("\n%d - %s", i+1, message)
		formatedMsg = append(formatedMsg, msg)
	}
	resultMsg := strings.Join(formatedMsg, "")
	return resultMsg
}
