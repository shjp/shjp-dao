package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	dao "github.com/shjp/shjp-dao"
	"github.com/shjp/shjp-dao/postgres"
)

func main() {
	envVars, err := godotenv.Read()
	if err != nil {
		panic(err)
	}

	addr := os.Getenv("SHJP_DB_HOST") + ":" + os.Getenv("SHJP_DB_PORT")
	user := os.Getenv("SHJP_DB_USER")
	dbName := os.Getenv("SHJP_DB_DATABASE")
	password := os.Getenv("SHJP_DB_PASSWORD")

	queueHost := envVars["QUEUE_URL"]
	queueUser := envVars["QUEUE_USER"]
	queueExchange := envVars["QUEUE_EXCHANGE"]

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

	/**
	 *	REST endpoints are used for query requests
	 */

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

	/**
	 * Subscrbies for mutation requests
	 */

	asyncService, err := dao.NewAsyncService(
		queueHost,
		queueUser,
		queueExchange,
		announcementService,
		eventService,
		groupService,
		userService,
		roleService)
	if err != nil {
		panic(err)
	}
	if err = asyncService.Listen(); err != nil {
		panic(err)
	}

	log.Println("Server listening on port 8200") //
	log.Fatal(http.ListenAndServe(":8200", r))
}
