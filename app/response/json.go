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
	W          http.ResponseWriter
	R          *http.Request
	Message    string
	StatusCode uint
	Data       interface{}
}

// JSONOK sends an http.StatusOK as the response together with the custom response `JSONResponse`.
func JSONOK(con Config) {

	// lets set the default value to ok
	if con.StatusCode <= 0 {
		con.StatusCode = http.StatusOK
	}

	render.JSON(con.W, con.R, &JSONResponse{
		StatusCode: con.StatusCode,
		Message:    con.Message,
		Success:    true,
		Data:       con.Data,
	})

	return
}

// JSONError handles the response for the client.
func JSONError(con Config) {
	render.JSON(con.W, con.R, &JSONResponse{
		StatusCode: con.StatusCode,
		Message:    con.Message,
		Success:    false,
		Data:       con.Data,
	})

	return
}
