package helpers

import "net/http"

// SendError writes data with specified status code to the response with status set to false.
//
// The response structure will be -
//
//	{
//	   "status": false,
//	   "message": message,
//	   "data": data
//	}
func SendError(w http.ResponseWriter, status int, message string, data interface{}) {
	SendJSON(w, status, map[string]any{
		"status":  false,
		"message": message,
		"data":    data,
	})
}
