package enrollment

import (
	"errors"
	"fmt"
)

var ErrorUserIDRequired = errors.New("user id is required")

var ErrorCourseIDRequired = errors.New("course id is required")

var ErrorStatusRequired = errors.New("status is required")

type ErrorEnrollmentNotFound struct {
	ID string
}

func (e ErrorEnrollmentNotFound) Error() string {
	return fmt.Sprintf("enrollment '%s' does not found", e.ID)
}

type ErrorInvalidStatus struct {
	Status string
}

func (e ErrorInvalidStatus) Error() string {
	return fmt.Sprintf("invalid '%s' status", e.Status)
}
