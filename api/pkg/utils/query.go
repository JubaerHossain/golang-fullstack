package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/entity"
)

func Paginate(req *http.Request, app *app.App, baseQuery, filterQuery string) (entity.Pagination, int, int, error) {
	ctx := req.Context()

	// Count total items with filters applied
	var totalItems int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s%s) AS filtered", baseQuery, filterQuery)
	if app.Config.DBType == "mysql" {
		if err := app.MDB.QueryRowContext(ctx, countQuery).Scan(&totalItems); err != nil {
			return entity.Pagination{}, 0, 0, err
		}
	} else {
		if err := app.DB.QueryRow(ctx, countQuery).Scan(&totalItems); err != nil {
			return entity.Pagination{}, 0, 0, err
		}
	}

	// Extract limit and offset from query parameters
	queryValues := req.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(queryValues.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	var nextPage, previousPage *int
	if page > 1 {
		prevPage := page - 1
		previousPage = &prevPage
	}
	if offset+limit < totalItems {
		nextPageValue := page + 1
		nextPage = &nextPageValue
	}

	// Calculate total pages
	totalPages := totalItems / limit
	if totalItems%limit != 0 {
		totalPages++
	}

	// Prepare pagination struct
	return entity.Pagination{
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     nextPage,
		PreviousPage: previousPage,
		FirstPage:    1,
		LastPage:     totalPages,
	}, limit, offset, nil
}
