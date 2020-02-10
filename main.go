package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
)

var sleeperProbes map[string]int

func init() {
	sleeperProbes = map[string]int{
		"1min":   60,
		"2min":   120,
		"4min":   240,
		"5min":   300,
		"15min":  900,
		"30min":  1800,
		"1hour":  3600,
		"2hour":  7200,
		"4hour":  14400,
		"8hour":  28000,
		"16hour": 57600,
		"1day":   86400,
		"2day":   172800,
		"4day":   345600,
		"1week":  604800,
		"2week":  1209600,
		"3week":  1814400,
		"4week":  2419200,
	}
}

func main() {
	var wg sync.WaitGroup

	fmt.Println("Starting", len(sleeperProbes), "probes to monitor DB connection health at the following intervals:")

	for probeName, _ := range sleeperProbes {
		fmt.Println("*", probeName)
	}

	for probeName, probeSec := range sleeperProbes {
		wg.Add(1)
		go sleeper(probeName, probeSec, &wg)
	}
	fmt.Println("\nestablishing connections...done")
	wg.Wait()
}

func getTime() string {
	return time.Now().Format("2006-01-02 03:04:05")
}

func sleeper(probeName string, probeSec int, wg *sync.WaitGroup) {
	var test int64
	var err error
	defer wg.Done()

	conn := pgConnect(probeName)
	defer conn.Close(context.Background())

	for {
		time.Sleep(time.Second * time.Duration(probeSec))
		err = conn.QueryRow(context.Background(), "select 1 as test").Scan(&test)
		if err != nil {
			//fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			fmt.Println(getTime(), probeName, "FAILURE", "(0.001s) --", err)
		} else {
			fmt.Println(getTime(), probeName, "OK", "(0.001s)")
		}
	}
}

func pgConnect(probeName string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	} else {
		//fmt.Println("Sucessfully open connection to database for ", probeName)
	}
	return conn
}
