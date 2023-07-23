package main

import (
  "database/sql"

  "encoding/json"

  "fmt"

  "log"

  "strings"

  "github.com/graphql-go/graphql"

  _ "github.com/denisenkom/go-mssqldb"
)

type Country struct {
  Abbr string
  Name string
}
type QueryCollection struct {
  Message string
  Query string
}

var countryType = graphql.NewObject(
    graphql.ObjectConfig{
        Name: "Country",
        Fields: graphql.Fields{
            "Abbr": &graphql.Field{
                Type: graphql.String,
            },
            "Name": &graphql.Field{
                Type: graphql.String,
            },
        },
    },
)
var mutationType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Mutation",
    Fields: graphql.Fields{
        "create": &graphql.Field{
            Type: countryType,
            Description: "Create a new Country",
            Args: graphql.FieldConfigArgument{
              "Abbr": &graphql.ArgumentConfig{
                  Type: graphql.NewNonNull(graphql.String),
              },
              "Name": &graphql.ArgumentConfig{
                  Type: graphql.NewNonNull(graphql.String),
              },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                country := Country{
                    Abbr: p.Args["Abbr"].(string),
                    Name: p.Args["Name"].(string),
                }
                countries = append(countries, country)
                return country, nil
            },
        },
    },
})

var countries []Country

func conn()  {

  const (
    host     = "db"
    user     = "sa"
    password = "z!oBx1ab"
    dbname   = "internet"
  )

  sqlserverInfo := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&connection+timeout=30",
     user, password, host, dbname)

  db, err := sql.Open("sqlserver", sqlserverInfo)

  if err != nil {
    panic(err)
  }

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("\nSuccessfully connected!")

  results, err00 := db.Query(`SELECT abbr, name FROM dbo.country ORDER BY abbr;`)
  if err00 != nil {
      panic(err00)
  }
  for results.Next() {
      var country Country
      err00 = results.Scan(&country.Abbr, &country.Name)
      if err00 != nil {
          fmt.Println(err00)
      }
      log.Println(country)
      countries = append(countries, country)
  }
  db.Close()
}

func nameContainsFilter(id string, arr []Country) []Country {

  results := make([]Country, 0)
  // Parse our tutorial array for the matching id
  for _, item := range arr {
      if strings.Contains(item.Name, id) {
          // return our tutorial
          results = append(results, item)
      }
  }
  return results
}
func queryResult(msg string, query string, schema graphql.Schema)  {
  params := graphql.Params{Schema: schema, RequestString: query}
  r := graphql.Do(params)
  if len(r.Errors) > 0 {
    log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
  }
  rJSON, _ := json.Marshal(r)
  fmt.Printf("%s:\n %s \n\n", msg, rJSON)
}
func main() {

  conn()

    // Schema
    fields := graphql.Fields{
        "Countries": &graphql.Field{
            Type: graphql.NewList(countryType),
            Args: graphql.FieldConfigArgument{
                "startWith": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
                "nameContains": &graphql.ArgumentConfig{
                    Type: graphql.String,
                },
            },
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
              startWith := make([]Country, 0)
              results := make([]Country, 0)
              var id, ok = p.Args["startWith"].(string)
                if ok {
                    // Parse our tutorial array for the matching id
                    for _, item := range countries {
                        if strings.HasPrefix(item.Abbr, id) {
                            // return our tutorial
                            startWith = append(startWith, item)
                        }
                    }
                }
                id, ok = p.Args["nameContains"].(string)
                  if ok {
                    // when startWith is used with this filter
                    if len(startWith) > 0 {
                        results = nameContainsFilter(id, startWith)
                    } else {
                      // when startWith is not used with this filter
                        results = nameContainsFilter(id, countries)
                    }
                  } else {
                    // when startWith is used without this filter
                      if len(startWith) > 0 {
                          results = startWith
                      } else {
                        // when neither filter is used
                          results = countries
                      }
                  }
                return results, nil
            },
        },
    }
    rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
    schemaConfig := graphql.SchemaConfig{
      Query: graphql.NewObject(rootQuery),
      Mutation: mutationType,
    }
    schema, err := graphql.NewSchema(schemaConfig)
    if err != nil {
        log.Fatalf("failed to create new schema, error: %v", err)
    }
    // Query
  	query := `
  		mutation {
  			create(Abbr: "eu", Name: "United States") {
  				Abbr
          Name
  			}
  		}
  	`
    queryResult("Mutation create", query, schema)

    querySlice := []QueryCollection {
      {Message:"Filtered by abbr: startWith", Query:`{Countries(startWith: "b"){Abbr Name}}`},
      {Message:"Filtered by name: nameContains", Query:`{Countries(nameContains: "z"){Abbr Name}}`},
      {Message:"Filtered by abbr: startWith and name: nameContains", Query:`{Countries(startWith: "b", nameContains: "z"){Abbr Name}}`},
      {Message:"No Filter", Query:`{Countries{Abbr Name}}`},
    }

    for _, q := range querySlice {
      queryResult(q.Message, q.Query, schema)
    }
}
