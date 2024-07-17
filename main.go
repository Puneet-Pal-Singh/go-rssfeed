package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main(){
	fmt.Println("Hello World")

	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("PORT is not bound in the environement")
	}

	fmt.Println("PORT:", portString)

	router := chi.NewRouter()

	// Create a new instance of the CORS handler
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"*"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: false,
        MaxAge:			  300, 
    })

    // Apply the CORS middleware
    router.Use(c.Handler)

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handleErr)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("server starting on port %v", portString)
	err := server.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}

}