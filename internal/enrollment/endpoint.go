package enrollment

import (
	"context"

	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Create endpoint.Endpoint
}

type CreateRequest struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		createRequest := request.(CreateRequest)

		if createRequest.UserID == "" {
			return nil, response.BadRequest(ErrorUserIDRequired.Error())
		}

		if createRequest.CourseID == "" {
			return nil, response.BadRequest(ErrorCourseIDRequired.Error())
		}

		enrollment, err := s.Create(
			createRequest.UserID,
			createRequest.CourseID,
		)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", enrollment, nil), nil
	}
}
