package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Puneet-Pal-Singh/go-rssfeed/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}

func main(){
	// feed, err := urlToFeed("https://www.wagslane.dev/index.xml")
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)
	
	fmt.Println("Hello World")

	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("PORT is not bound in the environement")
	}
	// fmt.Println("PORT:", portString)

	dbURL := os.Getenv("DB_URL")

	if dbURL == ""{
		log.Fatal("dbURL is not bound in the environement")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatal("Can't connect to database", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	log.Printf("started go routine")
	go startScraping(db, 10, time.Minute)
	log.Printf("go routine ended")
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
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreatefeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
	
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedsFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("server starting on port %v", portString)
	err = server.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}

}