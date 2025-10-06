package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/S3ergio31/curso-go-seccion-5-enrollment/internal/enrollment"
	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewEnrollmentHttpServer(endpoints enrollment.Endpoints) http.Handler {
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateEnrollment,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllEnrollment,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/enrollments", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateEnrollment,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	return router
}

func decodeCreateEnrollment(_ context.Context, r *http.Request) (any, error) {
	var request enrollment.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return request, nil
}

func decodeGetAllEnrollment(_ context.Context, r *http.Request) (any, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := enrollment.GetAllRequest{
		UserID:   v.Get("user_id"),
		CourseID: v.Get("course_id"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func decodeUpdateEnrollment(_ context.Context, r *http.Request) (any, error) {
	var request enrollment.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	path := mux.Vars(r)

	request.ID = path["id"]

	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, res any) error {
	r := res.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
