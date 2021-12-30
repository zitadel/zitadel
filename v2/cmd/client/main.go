package main

import (
	"context"
	"fmt"

	"github.com/caos/zitadel/pkg/grpc/management"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50002", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := management.NewManagementServiceClient(conn)

	languages, err := client.GetSupportedLanguages(context.Background(), &management.GetSupportedLanguagesRequest{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("languages: %v\n", languages.GetLanguages())
}
