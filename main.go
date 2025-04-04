package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/decoder"
)

func main() {
	content, err := os.ReadFile("dashboard.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file: %s\n", err)
		os.Exit(1)
	}
	handler := "405"
	modifiedContent := []byte(strings.ReplaceAll(string(content), "{{ handler }}", handler))

	dashboard, err := decoder.UnmarshalYAML(bytes.NewBuffer(modifiedContent))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse file: %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, "handler", handler)
	client := grabana.NewClient(&http.Client{}, "http://192.168.141.122:3000", grabana.WithAPIToken("eyJrIjoid1U5NjdGUDNIZ0xCdDN3YlNoa05jWWtuTEZQTk01UmEiLCJuIjoiQVBJIEpTT04iLCJpZCI6MX0="))

	// create the folder holding the dashboard for the service
	folder, err := client.FindOrCreateFolder(ctx, "Test Folder")
	if err != nil {
		fmt.Printf("Could not find or create folder: %s\n", err)
		os.Exit(1)
	}

	if _, err := client.UpsertDashboard(ctx, folder, dashboard); err != nil {
		fmt.Printf("Could not create dashboard: %s\n", err)
		os.Exit(1)
	}
}
