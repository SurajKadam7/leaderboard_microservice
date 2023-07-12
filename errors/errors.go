package youtubeerror

import "errors"

var ErrInvalidLimitValue = errors.New("invalid limit value passed")
var ErrEmptyVideoValuePassed = errors.New("video value shoud not be empty")

// db errors
var ErrVideoNotFound = errors.New("video no found")
var ErrNotAbleToIncrement = errors.New("not able to increase the view count")
var ErrNotAbleToDisplayTopViewed = errors.New("nt able to show the top viewed videos")