package main

import (
	"log"
	"os"

	"./hive"
)

func main() {
	client := new(hive.Client)
	client.Username = os.Getenv("EMAIL")
	client.Password = os.Getenv("PASS")
	err := client.GetData()
	if err != nil {
		log.Fatalf("Error getting data: %s", err)
	}
}
