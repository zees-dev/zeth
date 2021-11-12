package rest

const (
	HTTPBadRequest          = "Bad Request"
	HTTPInternalServerError = "Internal Server Error"
	HTTPNotFound            = "Not Found"
)

// HandlerError represents an error raised inside a HTTP handler
type HandlerError struct {
	StatusCode int
	Message    string
	Err        error
}

// // The serverError helper writes an error message and stack trace to the log,
// // then sends a custom response with specified http status code to the user.
// func serverError(w http.ResponseWriter, text string, statusCode int, err error) {
// 	http.Error(w, text, statusCode)
// }

// // The InternalServerError sends a 500 Internal Server Error response to the user.
// func InternalServerError(w http.ResponseWriter, err error) {
// 	log.Printf("%s\n", err.Error())
// 	serverError(w, "Internal Server Error", http.StatusInternalServerError, err)
// }

// // The BadRequest sends a 400 Bad Request Error response to the user.
// func BadRequest(w http.ResponseWriter, err error) {
// 	serverError(w, "Bad Request", http.StatusInternalServerError, err)
// }

// // The BadRequest sends a 404 Not Found Error response to the user.
// func NotFound(w http.ResponseWriter, err error) {
// 	serverError(w, "Not Found", http.StatusNotFound, err)
// }
