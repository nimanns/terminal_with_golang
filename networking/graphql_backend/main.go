package main

import (
    "log"
    "net/http"

    "github.com/graphql-go/graphql"
    "github.com/graphql-go/handler"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: "1", Name: "John Doe", Email: "john@example.com"},
    {ID: "2", Name: "Jane Smith", Email: "jane@example.com"},
}

func main() {
    userType := graphql.NewObject(graphql.ObjectConfig{
        Name: "User",
        Fields: graphql.Fields{
            "id": &graphql.Field{
                Type: graphql.String,
            },
            "name": &graphql.Field{
                Type: graphql.String,
            },
            "email": &graphql.Field{
                Type: graphql.String,
            },
        },
    })

    rootQuery := graphql.NewObject(graphql.ObjectConfig{
        Name: "RootQuery",
        Fields: graphql.Fields{
            "user": &graphql.Field{
                Type: userType,
                Args: graphql.FieldConfigArgument{
                    "id": &graphql.ArgumentConfig{
                        Type: graphql.String,
                    },
                },
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    id, ok := p.Args["id"].(string)
                    if ok {
                        for _, user := range users {
                            if user.ID == id {
                                return user, nil
                            }
                        }
                    }
                    return nil, nil
                },
            },
        },
    })

    schema, err := graphql.NewSchema(graphql.SchemaConfig{
        Query: rootQuery,
    })
    if err != nil {
        log.Fatalf("failed to create new schema, error: %v", err)
    }

    h := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
    })

    http.Handle("/graphql", h)
    log.Println("server is running on http://localhost:8080/graphql")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

