package enrollment

import (
	"errors"
)

var ErrorUserIDRequired = errors.New("user id is required")

var ErrorCourseIDRequired = errors.New("course id is required")
