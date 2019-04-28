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
	DB   *database.TwitterBankDatabase
	srv  *http.Server
	port int

	// Singleton Types
	apiTypes *GraphQLAPITypes
}

func NewTwitterBankApiServerFromDb(factomdHost string, factomdPort int, db *database.TwitterBankDatabase) (*TwitterBankApiServer, error) {
	s := new(TwitterBankApiServer)
	s.DB = db

	factom.SetFactomdServer(fmt.Sprintf("%s:%d", factomdHost, factomdPort))
	s.init()
	return s, nil
}

// ChangeListenPort can only be set before it is run
func (api *TwitterBankApiServer) ChangeListenPort(port int) {
	api.port = port
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

	log.Infof("Running on localhost:%d", api.port)
	api.srv = &http.Server{Addr: fmt.Sprintf(":%d", api.port)}
	http.Handle("/graphql", disableCors(h))
	err = api.srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
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
