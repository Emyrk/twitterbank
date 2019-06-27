package apiserver

import (
	"fmt"
	"log"

	"github.com/FactomProject/factom"
	"github.com/graphql-go/graphql"
)

func (api *TwitterBankApiServer) CreateSchema() (graphql.Schema, error) {
	// Schema
	fields := graphql.Fields{
		//"completed":            s.completedField(),
		//"proposal":             s.proposal(),
		//"allProposals":         s.allProposals(),
		//"eligibleList":         s.eligibleList(),
		//"eligibleVoters":       s.eligibleListVoters(),
		//"commit":               s.commit(),
		//"reveal":               s.reveal(),
		//"commits":              s.commits(),
		//"reveals":              s.reveals(),
		//"result":               s.result(),
		//"results":              s.results(),
		//"identityKeysAtHeight": s.identityKeysAtHeight(),
		//"proposalEntries":      s.proposalEntries(),
		"properties": api.Properties(),
		"user":       api.TwitterUser(),
		"tweet":      api.Tweet(),
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	return schema, err
}

func (api *TwitterBankApiServer) Properties() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:        "Properties",
			Description: "Various properties about the voting daemon",
			Fields: graphql.Fields{
				"syncedHeight": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return api.DB.FetchHighestDBInserted()
					},
				},
				"factomdProperties": &graphql.Field{
					Type: FactomdProperties,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						fdv, fdve, apv, apve, _, _, _, _ := factom.GetProperties()
						return []string{fdv, fdve, apv, apve}, nil
					},
				},
			}}),
		Description: "Returns various properties of the voting daemon",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return new(interface{}), nil
		},
	}
}

func (api *TwitterBankApiServer) FactomdProperties() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "FactomdProperties",
		Description: "Factomd Version",
		Fields: graphql.Fields{
			"factomdVersion": &graphql.Field{
				Type:        graphql.String,
				Description: "Version of factomd the scraper is talking to.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lst, ok := p.Source.([]string)
					if !ok {
						return nil, fmt.Errorf("Incorrect type supplied")
					}
					return lst[0], nil
				},
			},
			"factomdVersionError": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lst, ok := p.Source.([]string)
					if !ok {
						return nil, fmt.Errorf("Incorrect type supplied")
					}
					return lst[1], nil
				},
			},
			"factomdAPIVersion": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lst, ok := p.Source.([]string)
					if !ok {
						return nil, fmt.Errorf("Incorrect type supplied")
					}
					return lst[2], nil
				},
			},
			"factomdAPIVersionError": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lst, ok := p.Source.([]string)
					if !ok {
						return nil, fmt.Errorf("Incorrect type supplied")
					}
					return lst[3], nil
				},
			},
			"totalTrackedUsers": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return api.DB.FetchTotalNumberOfUsers()
				},
			},
			"totalNumberOfTweets": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return api.DB.FetchTotalNumberOfTweets()
				},
			},
			"totalNumberOfRecords": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return api.DB.FetchTotalNumberOfTweetRecords()
				},
			},
		}})
}
