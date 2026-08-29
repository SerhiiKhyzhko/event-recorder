package domain

import "errors"

var ErrInvalidDateRange = errors.New("start date must be before end date")