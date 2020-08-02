package hive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Username, Password string
	sessionID          string
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Redirect string
}

type payloadResponse struct {
	Location string
}

type DataPoint struct {
	Date        string
	Temperature float32
}

type historyResponse struct {
	Data []DataPoint `json:"data"`
}

type Product struct {
	Id string `json:"id"`
}

func (c *Client) getSession() error {
	client := &http.Client{}

	var creds = credentials{
		Username: c.Username,
		Password: c.Password,
		Redirect: "https://my.hivehome.com",
	}

	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", "https://sso.hivehome.com/auth/login?client=v3-web-prod", bytes.NewBuffer(data))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var payload payloadResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
	}

	parsedLocation, err := url.Parse(payload.Location)
	if err != nil {
		return err
	}

	m, err := url.ParseQuery(parsedLocation.Fragment)

	if err != nil {
		return err
	}

	c.sessionID = m.Get("id_token")

	return nil
}

func (c *Client) getProducts() ([]Product, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://beekeeper-uk.hivehome.com/1.0/products?after=", nil)
	req.Header.Set("authorization", c.sessionID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
	}
	return products, nil
}

func (c *Client) GetData(start time.Time, end time.Time) ([]DataPoint, error) {
	if c.sessionID == "" {
		err := c.getSession()
		if err != nil {
			log.Fatalf("Error getting session token: %s", err)
		}
	}

	client := &http.Client{}

	products, _ := c.getProducts()
	url := fmt.Sprintf("https://beekeeper-uk.hivehome.com/1.0/history/heating/%s?start=%d&end=%d&timeUnit=MINUTES&rate=30&operation=AVG", products[0].Id, start.Unix()*1000, end.Unix()*1000)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("authorization", c.sessionID)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var history historyResponse
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
	}

	if err != nil {
		return nil, err
	}

	return history.Data, nil
}
