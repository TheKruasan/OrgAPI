package main

import (
	"OrgAPI/internal/config"
	"OrgAPI/internal/database"
	"OrgAPI/internal/handler"
	"OrgAPI/internal/middleware"
	"OrgAPI/internal/repository"
	"OrgAPI/internal/service"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	// load env
	_ = godotenv.Load()

	// config
	cfg := config.Load()

	// database connection
	db, err := database.NewPostgres(cfg)

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal(err)
	}

	err = sqlDB.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected")

	// repositories
	departmentRepo := repository.NewDepartmentRepository(db)

	employeeRepo := repository.NewEmployeeRepository(db)

	// services
	departmentService := service.NewDepartmentService(
		departmentRepo,
		db,
	)

	employeeService := service.NewEmployeeService(
		employeeRepo,
		departmentRepo,
	)

	// handlers
	departmentHandler := handler.NewDepartmentHandler(
		departmentService,
	)

	employeeHandler := handler.NewEmployeeHandler(
		employeeService,
	)

	// router
	mux := http.NewServeMux()

	// department routes
	mux.HandleFunc(
		"POST /departments",
		departmentHandler.CreateDepartment,
	)

	mux.HandleFunc(
		"GET /departments/{id}",
		departmentHandler.GetDepartment,
	)

	mux.HandleFunc(
		"PATCH /departments/{id}",
		departmentHandler.UpdateDepartment,
	)

	mux.HandleFunc(
		"DELETE /departments/{id}",
		departmentHandler.DeleteDepartment,
	)

	// employee routes
	mux.HandleFunc(
		"POST /departments/{id}/employees",
		employeeHandler.CreateEmployee,
	)

	// server
	server := &http.Server{
		Addr: ":" + cfg.AppPort,
		Handler: middleware.Logging(
			mux,
		),
	}

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// run server
	go func() {

		log.Println(
			"server started on port",
			cfg.AppPort,
		)

		err := server.ListenAndServe()

		if err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(err)
		}
	}()

	// wait signal
	<-shutdown

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
