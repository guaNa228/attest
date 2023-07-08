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
	_ "github.com/lib/pq"
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

	err = apiCfg.createPrograms()

	if err != nil {
		fmt.Println("Error during programs initialization:", err)
	} else {
		fmt.Println("Programs are succesfully intialized!")
	}

	//go apiCfg.parsingResult()

	router := chi.NewRouter()

	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		},
	))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)

	v1Router.Get("/err", handlerErr)

	v1Router.Post("/user", apiCfg.middlewareAuth(apiCfg.handlerCreateUser, []string{}))

	v1Router.Post("/jwtCheck", apiCfg.middlewareAuth(apiCfg.handlerRoleByJWT, []string{"student", "teacher"}))

	v1Router.Post("/login", apiCfg.handlerLogin)

	v1Router.Post("/group", apiCfg.middlewareAuth(apiCfg.handlerCreateGroup, []string{}))
	v1Router.Delete("/group/{groupToDelete}", apiCfg.middlewareAuth(apiCfg.handlerDeleteGroup, []string{}))

	v1Router.Post("/class", apiCfg.middlewareAuth(apiCfg.handlerCreateClass, []string{}))
	v1Router.Delete("/class/{classToDeleteID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteClass, []string{}))

	v1Router.Post("/spawnAttestation", apiCfg.middlewareAuth(apiCfg.handleAttestationSpawn, []string{}))

	v1Router.Get("/attestation", apiCfg.middlewareAuth(apiCfg.handleAttestationGet, []string{"teacher", "student"}))
	v1Router.Post("/attestation", apiCfg.middlewareAuth(apiCfg.handleAttestationPost, []string{"teacher"}))

	v1Router.Post("/teacher", apiCfg.middlewareAuth(apiCfg.handlerCreateTeacher, []string{}))

	v1Router.Post("/student", apiCfg.middlewareAuth(apiCfg.handlerCreateStudent, []string{}))

	v1Router.Post("/students", apiCfg.middlewareAuth(apiCfg.uploadStudentsUpload, []string{}))

	v1Router.Post("/parsing", apiCfg.middlewareAuth(apiCfg.handleParsing, []string{}))

	v1Router.Post("/mails_parsing", apiCfg.middlewareAuth(apiCfg.handleEmailParsing, []string{}))

	v1Router.Route("/", func(ws chi.Router) {
		ws.Get("/ws", wsHandler)
	})

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
