// main
// Simple App with mongo db

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	setroutes(router)
	port := ":8080"
	log.Printf("API server listening at %s", port)
	log.Fatal(http.ListenAndServe(port, router))
}

// Setup server routes
// Setting server datapoints handllers routes
func setroutes(router *mux.Router) {
	router.HandleFunc("/health", health).Methods("GET")
	router.HandleFunc("/", app).Methods("GET")
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to Test API service.\n Health: Ok")
}

func app(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI(getMongoConn()))
	if err != nil {
		fmt.Fprintf(w, "Error while creating to mongo client, err: %s", err)
		return
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Fprintf(w, "Error while connecting to mongo db, err: %s", err)
		return
	}
	defer client.Disconnect(ctx)

	// List databases
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		fmt.Fprintf(w, "Error while getting collection list, err: %s", err)
		return
	}

	fmt.Fprintf(w, `<style>
	div {
		max-width: 500px;
		margin: auto;
		align: center;
	  }
	table,tbody {
		border: 1px solid black;
		display: table;
		border-collapse: collapse;
		border-spacing: 2px;
		border-color: gray;
		font-size: 16;
	  } </style> <div><br/><h2>Welcome to demo APP</h2>
	  				<br />Mongo  db Connection: OK <br /><br />
					<table><tr><th>Collections</th></tr>
					<tr><td>%s<td /><tr />
					</table></div>`, strings.Join(databases, ", "))
}

func getMongoConn() string {
	host := os.Getenv("MONGO_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("MONGO_PORT")
	if port == "" {
		port = "27017"
	}
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}
