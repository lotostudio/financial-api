package v1

import "errors"

var (
	errDateFiltersInvalid = errors.New("date filters are invalid. check 'dateFrom' and 'dateTo' params")
)
