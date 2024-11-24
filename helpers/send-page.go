package helpers

import "net/http"

// Page defines the data structure send for pagination response's data field.
type Page struct {
	// Items are the list items
	Items interface{} `json:"items"`

	// ItemsPerPage is the number of items per page
	ItemsPerPage int `json:"itemsPerPage"`

	// PageNumber is the current page number
	PageNumber int `json:"pageNumber"`

	// TotalItems is the total number of items found
	TotalItems int `json:"totalItems"`

	// TotalPages is the total number of pages available
	TotalPages int `json:"totalPages"`
}

// SendPage writes a pagination response in the pre-defined response structure.
//
// Pagination response structure -
//
//	{
//	   "status": true,
//	   "message": "Success",
//	   "data": {
//	     "pageNumber": number,
//	     "itemsPerPage": number,
//	     "totalItems": number,
//	     "totalPages": number,
//	     "items": []any
//	   }
//	}
func SendPage(w http.ResponseWriter, page Page) {
	SendData(w, page)
}
