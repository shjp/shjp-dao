package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

	dao "github.com/shjp/shjp-dao"
	"github.com/shjp/shjp-dao/postgres"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	reqBlob, err := json.Marshal(request)
	if err != nil {
		log.Println("Marshalling request failed:", err)
	}
	log.Println("Request object ---------------------------------------------------")
	log.Println(string(reqBlob))
	log.Println("------------------------------------------------------------------")

	// authToken, ok := request.Headers["auth-token"]
	// // For time being, simply log and pass an empty string when auth token is not found
	// if !ok {
	// 	log.Println("Auth token not found")
	// }

	addr := os.Getenv("SHJP_DB_HOST") + ":" + os.Getenv("SHJP_DB_PORT")
	user := os.Getenv("SHJP_DB_USER")
	dbName := os.Getenv("SHJP_DB_DATABASE")
	password := os.Getenv("SHJP_DB_PASSWORD")

	log.Print("Initializing DB client...")
	db := postgres.Init(&pg.Options{
		Addr:     addr,
		Password: password,
		User:     user,
		Database: dbName,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	// Log queries
	db.AddQueryHook(postgres.Logger{})

	announcementService := dao.NewModelService(&postgres.AnnouncementQueryStrategy{DB: db})
	eventService := dao.NewModelService(&postgres.EventQueryStrategy{DB: db})
	groupService := dao.NewModelService(&postgres.GroupQueryStrategy{DB: db})
	userService := dao.NewModelService(&postgres.UserQueryStrategy{DB: db})
	roleService := dao.NewModelService(&postgres.RoleQueryStrategy{DB: db})

	r := mux.NewRouter()
	r.Path("/announcements").HandlerFunc(announcementService.HandleGetAll)
	r.Path("/announcements/search").HandlerFunc(announcementService.HandleSearch)
	r.Path("/announcements/{id}").HandlerFunc(announcementService.HandleGetOne)
	r.Path("/events").HandlerFunc(eventService.HandleGetAll)
	r.Path("/events/search").HandlerFunc(eventService.HandleSearch)
	r.Path("/events/{id}").HandlerFunc(eventService.HandleGetOne)
	r.Path("/groups").HandlerFunc(groupService.HandleGetAll)
	r.Path("/groups/search").HandlerFunc(groupService.HandleSearch)
	r.Path("/groups/{id}").HandlerFunc(groupService.HandleGetOne)
	r.Path("/users").HandlerFunc(userService.HandleGetAll)
	r.Path("/users/search").HandlerFunc(userService.HandleSearch)
	r.Path("/users/{id}").HandlerFunc(userService.HandleGetOne)
	r.Path("/roles").HandlerFunc(roleService.HandleGetAll)
	r.Path("/roles/search").HandlerFunc(roleService.HandleSearch)
	r.Path("/roles/{id}").HandlerFunc(roleService.HandleGetOne)

	return handleLambdaEvent(r.ServeHTTP, request)
}
