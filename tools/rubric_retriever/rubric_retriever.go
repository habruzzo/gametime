package rubric_retriever

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
)

func main() {
	keyFile := os.Getenv("GOOGLE_CRED_LOCATION")
	client, err := sheets.NewService(context.Background(), option.WithCredentialsFile(keyFile))
	if err != nil {
		panic(fmt.Sprintf("error getting client! %s", err.Error()))
	}
	s := client.Spreadsheets.Get("1PZY2BLdhYFiWJBJVijD3MPsErVEdA-9aXv0CzN_-Zig")
	s.Do()
}
