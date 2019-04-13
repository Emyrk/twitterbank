package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Emyrk/twitterbank/database"
	"github.com/FactomProject/factom"
	"github.com/graphql-go/handler"
	log "github.com/sirupsen/logrus"
)

type TwitterBankApiServer struct {
	DB  *database.TwitterBankDatabase
	srv *http.Server
}

func NewTwitterBankApiServerFromDb(factomdHost string, factomdPort int, db *database.TwitterBankDatabase) (*TwitterBankApiServer, error) {
	s := new(TwitterBankApiServer)
	s.DB = db

	factom.SetFactomdServer(fmt.Sprintf("%s:%d", factomdHost, factomdPort))
	return s, nil
}

func (api *TwitterBankApiServer) RunDaemon() {
	schema, err := api.CreateSchema()
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:       &schema,
		Pretty:       true,
		GraphiQL:     false,
		Playground:   true,
		RootObjectFn: api.rootObjF,
	})

	port := 8080
	log.Infof("Running on localhost:%d", port)
	api.srv = &http.Server{Addr: fmt.Sprintf(":%d", port)}
	http.Handle("/graphql", disableCors(h))
	api.srv.ListenAndServe()
}

func (api *TwitterBankApiServer) Close() error {
	return api.srv.Close()
}

func (api *TwitterBankApiServer) rootObjF(ctx context.Context, r *http.Request) map[string]interface{} {
	return map[string]interface{}{"apiserver": api}
}

// disableCors from: https://github.com/graphql-go/graphql/issues/290
func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		h.ServeHTTP(w, r)
	})
}
