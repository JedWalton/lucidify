package lucidifychat

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

func GetSchema() {
	cfg := weaviate.Config{
		Host:   "localhost:8090",
		Scheme: "http",
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", schema)
}
