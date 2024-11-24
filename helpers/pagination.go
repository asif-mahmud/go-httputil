package helpers

import (
	"math"
	"net/http"
	"strconv"
)

// FilterParams is collection of query parameters found in the search queries of the incoming request.
// This will contain only the parameters which are not defined by the pagination constants like
// pageNumber, itemsPerPage, search etc.
type FilterParams map[string][]string

// PaginationParams container for pagination parameters collected from url search queries.
type PaginationParams struct {
	// Page is the requested page number
	Page int

	// Limit is the number of items requested per page
	Limit int

	// Search is the common search term query
	Search string

	// SortBy is the name of the property which should be used for sorting the list
	SortBy string

	// SortOrder is the requested sort order.
	// It's value can be application specific, for example asc, desc could be
	// used to sort in ascending and descending order respectively, or any other
	// value defined by the application.
	SortOrder string

	// Filters are the additional filter parameters requested.
	// These are also application specific application may define any
	// filter query for any REST API, for example, <url>?isActive=true
	// defines a filter isActive with value true.
	Filters FilterParams
}

const (
	// MaxPageSize limits the maximum possible page size.
	MaxPageSize = 100

	// DefaultPageSize sets the default page size.
	DefaultPageSize = 10

	// PageKey is the search query key for PaginationParams.Page.
	PageKey = "pageNumber"

	// LimitKey is the search query key for PaginationParams.Limit.
	LimitKey = "itemsPerPage"

	// SearchKey is the search query key for PaginationParams.Search.
	SearchKey = "search"

	// SortByKey is the search query key for the PaginationParams.SortBy.
	SortByKey = "sortBy"

	// SortOrderKey is the search query key for the PaginationParams.SortOrder.
	SortOrderKey = "sortOrder"
)

func parsePage(r *http.Request) int {
	pageStr := r.URL.Query().Get(PageKey)
	page, _ := strconv.ParseInt(pageStr, 10, 32)
	page = int64(math.Max(1.0, float64(page)))
	return int(page)
}

func parseLimit(r *http.Request) int {
	limitStr := r.URL.Query().Get(LimitKey)
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	limit = int64(math.Max(0.0, math.Min(MaxPageSize, float64(limit))))
	return int(limit)
}

// CountTotalPages counts total number of pages from total item count and per page item count.
func CountTotalPages(limit, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / math.Max(1.0, float64(limit))))
}

// GetPaginationParams parses the request's search query and constructs PaginationParams instance.
// There are a couple of pre-defined keys for pre-defined puposes, check the constants defined in
// this package. Any other key which is not a pre-defined key will be considered as an
// application specific filter query and added to the PaginationParams.FilterParams.
func GetPaginationParams(r *http.Request, defaultSortBy, defaultSortOrder string) PaginationParams {
	params := PaginationParams{
		Page:      1,
		Limit:     DefaultPageSize,
		Search:    "",
		SortBy:    defaultSortBy,
		SortOrder: defaultSortOrder,
		Filters:   FilterParams{},
	}

	for k, v := range r.URL.Query() {
		switch k {
		case PageKey:
			// parse page number
			params.Page = parsePage(r)

		case LimitKey:
			// parse limit
			params.Limit = parseLimit(r)

		case SearchKey:
			// parse search term
			params.Search = r.URL.Query().Get(SearchKey)

		case SortByKey:
			// parse sort by
			params.SortBy = r.URL.Query().Get(SortByKey)

		case SortOrderKey:
			// parse sort order
			params.SortOrder = r.URL.Query().Get(SortOrderKey)

		default:
			// any other filter parameter
			params.Filters[k] = v
		}
	}

	return params
}

// GetSortingData parses sorting related parameters from url search query.
// This should not be directly used as it only calls GetPaginationParams
// under the hood.
func GetSortingData(
	r *http.Request,
	defaultSortBy, defaultSortOrder string,
) (sortBy, sortOrder string) {
	params := GetPaginationParams(r, defaultSortBy, defaultSortOrder)
	return params.SortBy, params.SortOrder
}
