package hive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	Username, Password string
	sessionID          string
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Caller   string `json:"caller"`
}

type payload struct {
	Sessions []credentials `json:"sessions"`
}

type sessions struct {
	SessionID string `json:"sessionId"`
}

type payloadResponse struct {
	Sessions []sessions `json:"sessions"`
}

func (c *Client) GetSession() error {
	// fmt.Println("Getting new session token")
	client := &http.Client{}

	var creds = credentials{
		Username: c.Username,
		Password: c.Password,
		Caller:   "WEB",
	}

	var p = payload{
		Sessions: []credentials{creds},
	}

	data, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("JSON Marshal error: %s", err)
	}

	u := bytes.NewReader(data)

	req, _ := http.NewRequest("POST", "https://api-prod.bgchprod.info:443/omnia/auth/sessions", u)
	req.Header.Set("Content-Type", "application/vnd.alertme.zoo-6.1+json")
	req.Header.Set("Accept", "application/vnd.alertme.zoo-6.1+json")
	req.Header.Set("X-Omnia-Client", "Hive Web Dashboard")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var sessions payloadResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		return err
	}
	c.sessionID = sessions.Sessions[0].SessionID
	return nil
}

func (c *Client) GetData() error {
	if c.sessionID == "" {
		err := c.GetSession()
		if err != nil {
			log.Fatalf("Error getting session token: %s", err)
		}
	}
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api-prod.bgchprod.info:443/omnia/nodes", nil)
	req.Header.Set("Content-Type", "application/vnd.alertme.zoo-6.1+json")
	req.Header.Set("Accept", "application/vnd.alertme.zoo-6.1+json")
	req.Header.Set("X-Omnia-Client", "Hive Web Dashboard")
	req.Header.Set("X-Omnia-Access-Token", c.sessionID)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error 1")
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error 2")
		return err
	}
	fmt.Printf("%s\n", body)
	return nil
}
