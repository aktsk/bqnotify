package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aktsk/bqnotify/lib/runner"
)

// BqNotify is a main function to run bqnotify on Cloud Functions
func BqNotify() error {
	log.Printf("GCP_PROJECT: %s\n", os.Getenv("GCP_PROJECT"))
	err := runner.Run("./config.yaml")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	server := http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var jb map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&jb); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		log.Printf("Received JSON: %v\n", jb)
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("ListenAndServe: %v", err)
		}
	}()
	os.Setenv("SLACK_WEBHOOK_URL", "http://localhost:8080")
	if err := BqNotify(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
