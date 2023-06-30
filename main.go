package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	db "github.com/guaNa228/attest/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {

	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the enviroment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the enviroment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	db := db.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	fmt.Println("Succesfully connected to database")

	router := chi.NewRouter()

	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		},
	))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)

	v1Router.Get("/err", handlerErr)

	v1Router.Post("/user", apiCfg.middlewareAuth(apiCfg.handlerCreateUser, []string{}))

	v1Router.Post("/login", apiCfg.handlerLogin)

	v1Router.Get("/test", apiCfg.middlewareAuth(apiCfg.handlerGetUser, []string{"teacher", "student"}))

	v1Router.Post("/group", apiCfg.middlewareAuth(apiCfg.handlerCreateGroup, []string{}))
	v1Router.Delete("/group/{groupToDelete}", apiCfg.middlewareAuth(apiCfg.handlerDeleteGroup, []string{}))

	v1Router.Post("/class", apiCfg.middlewareAuth(apiCfg.handlerCreateClass, []string{}))
	v1Router.Delete("/class/{classToDeleteID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteClass, []string{}))

	v1Router.Post("/semesterActivity", apiCfg.middlewareAuth(apiCfg.handlerCreateSemesterActivity, []string{}))
	v1Router.Post("/semesterActivity/{semesterActivityToUpdateID}", apiCfg.middlewareAuth(apiCfg.handlerUpdateSemesterActivity, []string{}))
	v1Router.Delete("/semesterActivity/{semesterActivityToDeleteID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteSemesterActivity, []string{}))

	v1Router.Post("/attestation", apiCfg.middlewareAuth(apiCfg.handleAttestationSpawn, []string{}))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v", port)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
