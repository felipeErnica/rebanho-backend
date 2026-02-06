package util

import (
	"net/http"
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
