package helpers

import "net/http"

// SendData writes JSON marshalled data to the response in pre-defined structure.
//
// The response structure will be -
//
//	{
//	  "status": true,
//	  "message": "Success",
//	  "data": data
//	}
func SendData(w http.ResponseWriter, data interface{}) {
	sendJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Success",
		"data":    data,
	})
}
