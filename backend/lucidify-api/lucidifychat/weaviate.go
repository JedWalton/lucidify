package lucidifychat

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func CreateWeaviateClass() {
	cfg := weaviate.Config{
		Host:   "docker-weaviate-1:8080",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	// schema, err := client.Schema().ClassGetter().WithClassName("Article").Do(context.Background())
	// if err != nil {
	// 	panic(err)
	// }
	//
	// log.Printf("Schema:  %#v", schema)
}
