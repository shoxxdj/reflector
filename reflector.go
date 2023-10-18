package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Elem struct {
	Url  string
	Tags []string
}

func main() {
	urlFlag := flag.String("u", "", "URL to test")
	modeFlag := flag.String("m", "sniper", "Mode of fuzzing\n\t\tsniper= one param a time\n\t\tbram= all parameters same time with different values\n\t\tanything else : Same pattern in all the parameters")
	verboseFlag := flag.Bool("v", false, "Verbose mode")
	colorFlag := flag.Bool("c", false, "Color mode")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "            __ _           _             \n           / _| |         | |            \n  _ __ ___| |_| | ___  ___| |_ ___  _ __ \n | '__/ _ \\  _| |/ _ \\/ __| __/ _ \\| '__|\n | | |  __/ | | |  __/ (__| || (_) | |   \n |_|  \\___|_| |_|\\___|\\___|\\__\\___/|_|                                            \n                                         \n")
		fmt.Fprintf(os.Stderr, "Reflector : a binary to find xss based on url params reflections. v:0.1\n")

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "\t-%v: %v\n", f.Name, f.Usage) // f.Name, f.Value
		})
	}

	// Parse the command-line flags.
	flag.Parse()

	// Initialize a slice to store URLs.
	var urls []string

	// Check if the URL was provided as a flag, and add it to the slice.
	if *urlFlag != "" {
		urls = append(urls, *urlFlag)
	}

	// Read URLs from stdin and add them to the slice.
	var input string

	if len(*urlFlag) == 0 {
		for {
			_, err := fmt.Scanln(&input)
			if err != nil {
				break
			}
			if input != "" {
				urls = append(urls, input)
			}
		}
	}

	// Parse and process each URL in the slice.
	for _, inputURL := range urls {
		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			fmt.Printf("Error parsing URL %s: %v\n", inputURL, err)
			continue
		}

		// fmt.Println("URL:", inputURL)
		// fmt.Println("Scheme:", parsedURL.Scheme)
		// fmt.Println("Host:", parsedURL.Host)
		// fmt.Println("Path:", parsedURL.Path)
		// fmt.Println("Raw Query:", parsedURL.RawQuery)

		queryValues, err := url.ParseQuery(parsedURL.RawQuery)

		// Create a new URL with one query parameter replaced by a random value.
		newURL := *parsedURL
		originalQuery := queryValues.Encode()

		var jobsToDo []Elem

		if *modeFlag == "sniper" {
			// Iterate through the query parameters, replacing one at a time.
			for key := range queryValues {
				// Generate a random value.
				randomValue := generateRandomValue()

				// Temporarily replace the value for the current query parameter.
				originalValue := queryValues.Get(key)
				queryValues.Set(key, randomValue)

				// Update the query string in the new URL.
				newURL.RawQuery = queryValues.Encode()
				//fmt.Println(newURL.String())

				var tags []string
				tags = append(tags, randomValue)

				newObj := Elem{Url: newURL.String(), Tags: tags}
				jobsToDo = append(jobsToDo, newObj)

				// Restore the original value for the next iteration.
				queryValues.Set(key, originalValue)
				newURL.RawQuery = originalQuery
			}
		} else if *modeFlag == "bram" {
			randomValue := generateRandomValue()
			for key := range queryValues {
				queryValues.Set(key, randomValue)
			}
			newURL.RawQuery = queryValues.Encode()
			var tags []string
			tags = append(tags, randomValue)
			newObj := Elem{Url: newURL.String(), Tags: tags}

			jobsToDo = append(jobsToDo, newObj)
		} else {
			var tags []string
			for key := range queryValues {
				// Generate a random value.
				randomValue := generateRandomValue()
				queryValues.Set(key, randomValue)
				newURL.RawQuery = queryValues.Encode()
				tags = append(tags, randomValue)
			}
			newObj := Elem{Url: newURL.String(), Tags: tags}
			jobsToDo = append(jobsToDo, newObj)
		}

		//fmt.Println(jobsToDo)

		for _, elem := range jobsToDo {
			resp, err := http.Get(elem.Url)
			if err != nil {
				log.Fatalln(err)
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			sb := string(body)

			for _, tag := range elem.Tags {
				if strings.Contains(sb, tag) {
					if *colorFlag {
						color.Green(elem.Url)
					} else {
						fmt.Println(elem.Url)
					}

				} else if *verboseFlag {
					if *colorFlag {
						color.Red("Not Found")
					} else {
						fmt.Println("Not Found")
					}
				}
			}

		}

		//Generate object "url"=>"random value" & push dans un array & pour chaque elem du array send les req

		//Mode full fuzz => tous les params d'un coup
		//Mode onetoone => Un param à la fois
	}
}

func generateRandomValue() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomValue := make([]byte, 10)
	for i := range randomValue {
		randomValue[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomValue)
}
