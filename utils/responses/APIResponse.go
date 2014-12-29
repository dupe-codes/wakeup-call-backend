package APIResponses

import (
    "encoding/json"
    "net/http"
    "fmt"
)

type Response struct {
	Status int
	Data   interface{}
	Error  error
}

func SendErrorResponse(errorMsg error, status int, res http.ResponseWriter) {
    resContent := &Response{Status: status, Error: errorMsg}
    response, err := json.MarshalIndent(resContent, "", "  ")
    if err != nil {
        http.Error(res, "Error preparing response", http.StatusInternalServerError)
        return
    }
    http.Error(res, string(response), status)
    return
}

func SendSuccessResponse(data string, res http.ResponseWriter) {
    resContent := &Response{Status: http.StatusOK, Data: data}
    response, err := json.MarshalIndent(resContent, "", "  ")
    if err != nil {
        http.Error(res, "Error preparing response", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(res, string(response))
    return
}
