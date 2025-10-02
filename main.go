package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/S3ergio31/curso-go-seccion-5-enrollment/internal/enrollment"
	"github.com/S3ergio31/curso-go-seccion-5-enrollment/pkg/bootstrap"
	"github.com/S3ergio31/curso-go-seccion-5-enrollment/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := bootstrap.DBConnection()
	logger := bootstrap.InitLogger()

	if err != nil {
		logger.Fatalln(err)
	}

	enrollmentRepository := enrollment.NewRepository(logger, db)
	enrollmentService := enrollment.NewService(enrollmentRepository, logger)
	enrollmentEndpoints := enrollment.MakeEndpoints(enrollmentService)

	address := fmt.Sprintf("%s:%s", os.Getenv("APP_URL"), os.Getenv("APP_PORT"))
	server := &http.Server{
		Handler:      handler.NewEnrollmentHttpServer(enrollmentEndpoints),
		Addr:         address,
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	errCh := make(chan error)
	go func() {
		logger.Println("listen in ", address)
		errCh <- server.ListenAndServe()
	}()

	err = <-errCh

	if err != nil {
		logger.Fatal(err)
	}
}
