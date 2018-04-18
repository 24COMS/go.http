package server

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// GetLimitAndOffsetParams parses offset anf limit from url values
func GetLimitAndOffsetParams(values url.Values) (offset uint64, limit uint64, err error) {
	// TODO: in some time we would have to create common solution to avoid copy/paste
	rawOffset := values.Get("offset")
	rawLimit := values.Get("limit")

	if len(rawOffset) > 0 {
		offset, err = strconv.ParseUint(rawOffset, 10, 64)
		if err != nil {
			err = errors.Wrap(err, "failed to parse offset")
			return
		}
	}

	if len(rawLimit) > 0 {
		limit, err = strconv.ParseUint(rawLimit, 10, 64)
		if err != nil {
			err = errors.Wrap(err, "failed to parse limit")
			return
		}
	}
	return
}
