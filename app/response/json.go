package response

import (
	"net/http"

	"github.com/go-chi/render"
)

// JSONResponse custom json response
type JSONResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
}

// Config sets the different response configuration when returning a JSON responses.
type Config struct {
	Message    string
	StatusCode uint
	Data       interface{}
}

// Configure configures the response by a given message, statusCode, data.
func Configure(message string, statusCode uint, data interface{}) Config {
	return Config{
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}

// JSONOK sends an http.StatusOK as the response together with the custom response `JSONResponse`.
func JSONOK(w http.ResponseWriter, r *http.Request, con Config) {

	// lets set the default value to ok
	if con.StatusCode <= 0 {
		con.StatusCode = http.StatusOK
	}

	render.JSON(w, r, &JSONResponse{
		StatusCode: con.StatusCode,
		Message:    con.Message,
		Success:    true,
		Data:       con.Data,
	})

	return
}

// JSONError handles the response for the client.
func JSONError(w http.ResponseWriter, r *http.Request, con Config) {
	render.JSON(w, r, &JSONResponse{
		StatusCode: con.StatusCode,
		Message:    con.Message,
		Success:    false,
		Data:       con.Data,
	})

	return
}
