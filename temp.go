package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"./hive"
)

func main() {
	client := new(hive.Client)
	client.Username = os.Getenv("EMAIL")
	client.Password = os.Getenv("PASS")

	end := time.Now()
	start := end.AddDate(0, 0, -10)
	data, err := client.GetData(start, end)
	if err != nil {
		log.Fatalf("Error getting data: %s", err)
	}

	const format = "%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Time", "Temperature")
	fmt.Fprintf(tw, format, "----", "-----------")
	for _, v := range data {
		ts, _ := strconv.ParseInt(v.Date, 0, 64)
		fmt.Fprintf(tw, format, time.Unix(ts/1000, 0), v.Temperature)
	}
	tw.Flush()
}
