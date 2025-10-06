package enrollment

import (
	"context"
	"errors"
	"os"

	"github.com/S3ergio31/curso-go-seccion-5-meta/meta"
	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Create endpoint.Endpoint
	GetAll endpoint.Endpoint
	Update endpoint.Endpoint
}

type CreateRequest struct {
	UserID   string `json:"user_id"`
	CourseID string `json:"course_id"`
}

type GetAllRequest struct {
	UserID   string
	CourseID string
	Limit    int
	Page     int
}

type UpdateRequest struct {
	ID     string
	Status *string `json:"status"`
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
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

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetAllRequest)
		filters := Filters{
			UserID:   req.UserID,
			CourseID: req.CourseID,
		}

		count, err := s.Count(filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, os.Getenv("PAGINATOR_LIMIT_DEFAULT"))

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		enrollments, err := s.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		return response.Ok("success", enrollments, nil), nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		updateRequest := request.(UpdateRequest)

		if updateRequest.Status != nil && *updateRequest.Status == "" {
			return nil, response.BadRequest(ErrorStatusRequired.Error())
		}

		err := s.Update(
			updateRequest.ID,
			updateRequest.Status,
		)

		if errors.As(err, &ErrorEnrollmentNotFound{}) {
			return nil, response.NotFound(err.Error())
		}

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Ok("success", nil, nil), nil
	}
}
