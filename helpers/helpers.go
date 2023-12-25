package helpers

import (
	"encoding/json"
	"errors"
	"github.com/akashc777/csvToPdf/services"
	"log"
	"net/http"
	"os"
)

type ContextKey string

const (
	UserEmailContextKey ContextKey = "userInfo"
)

type Envolope map[string]interface{}

type Message struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var MessageLogs = &Message{
	InfoLog:  infoLog,
	ErrorLog: errorLog,
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxByte := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != nil {
		return errors.New("Body must have only a single json object")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		MessageLogs.ErrorLog.Printf(
			"helpers/WriteJSON Failed to unmarshal data err: %+v",
			err)
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}

	}

	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		MessageLogs.ErrorLog.Printf(
			"helpers/WriteJSON Failed to write to response header err: %+v",
			err)
	}

}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload services.JsonResponse
	payload.Error = true
	payload.Message = err.Error()
	WriteJSON(w, statusCode, payload)
}
