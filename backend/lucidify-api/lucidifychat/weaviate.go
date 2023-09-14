package lucidifychat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

func CreateWeaviateClass() {
	cfg := weaviate.Config{
		Host:   "docker-weaviate-1:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Println(err)
	}

	classObj := &models.Class{
		Class:       "Documents",
		Description: "A written text document",
		Properties: []*models.Property{
			{
				DataType:    []string{"string"},
				Description: "Title of the document",
				Name:        "title",
			},
			{
				DataType:    []string{"text"},
				Description: "The content of the document",
				Name:        "content",
			},
		},
	}

	err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		log.Printf("Documents Class already exists, not updating schema: %#v", err)
	}

	schema, err := client.Schema().ClassGetter().WithClassName("Documents").Do(context.Background())
	if err != nil {
		panic(err)
	}

	log.Printf("Schema:  %#v", schema)
}

func CreateDataObjects() {
	cfg := weaviate.Config{
		Host:   "docker-weaviate-1:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	dataSchema := map[string]interface{}{
		"name":    "Jodi Kantor",
		"content": "Hello, world!",
	}

	created, err := client.Data().Creator().
		WithClassName("Documents").
		WithID("36ddd591-2dee-4e7e-a3cc-eb86d30a4303").
		WithProperties(dataSchema).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL). // default QUORUM
		Do(context.Background())

	if err != nil {
		log.Println(err)
	}
	log.Printf("%v", created)
}

func ListObjectsInWeaviateClass() {
	cfg := weaviate.Config{
		Host:   "docker-weaviate-1:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	meta := graphql.Field{
		Name: "meta", Fields: []graphql.Field{
			{Name: "count"},
		},
	}

	result, err := client.GraphQL().Aggregate().
		WithClassName("Documents").
		WithFields(meta).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("ListObjectsInWeaviateClass: %v", result)
}

func GenerativeSearch() {
	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	cfg := weaviate.Config{
		Host:   "docker-weaviate-1:8080",
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": OPENAI_API_KEY, // Replace with your API key
		},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	name := graphql.Field{Name: "name"}

	concepts := []string{"magazine or newspaper about finance"}
	nearText := client.GraphQL().NearTextArgBuilder().
		WithConcepts(concepts)

	gs := graphql.NewGenerativeSearch().GroupedResult("Explain why these magazines or newspapers are about finance")

	result, err := client.GraphQL().Get().
		WithClassName("Documents").
		WithFields(name).
		WithGenerativeSearch(gs).
		WithNearText(nearText).
		WithLimit(5).
		Do(ctx)

	if err != nil {
		panic(err)
	}
	log.Printf("%v", result)

	type Response struct {
		Data struct {
			Get struct {
				Publication []struct {
					Additional struct {
						Generate struct {
							Error         interface{} `json:"error"`
							GroupedResult string      `json:"groupedResult"`
						} `json:"_additional"`
					} `json:"generate"`
					Name string `json:"name"`
				} `json:"Publication"`
			} `json:"Get"`
		} `json:"data"`
	}
	responseJSON, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	parsedResult := &Response{}
	err = json.Unmarshal(responseJSON, parsedResult)
	if err != nil {
		panic(err)
	}

	for _, publication := range parsedResult.Data.Get.Publication {
		fmt.Println("Name:", publication.Name)
		if publication.Additional.Generate.GroupedResult != "" {
			fmt.Println("Generated Explanation:", publication.Additional.Generate.GroupedResult)
		}
		fmt.Println("-----------")
	}
}
