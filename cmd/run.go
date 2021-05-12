package cmd

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const baseUrl = "http://localhost:8080"

var wg sync.WaitGroup

func Run() {
	log.Println("requests started")

	for x := 1; x <= 10; x++ {
		wg.Add(1)

		go func(step int) {
			defer wg.Done()

			time.Sleep(time.Second * time.Duration(step))

			c := http.Client{}
			res, err := c.Get(baseUrl)
			if err != nil {
				log.Printf("error occurred at step %d. err: %s\n", step, err)

				return
			}

			log.Printf("Request step %d status code %d\n", step, res.StatusCode)
		}(x)
	}

	wg.Wait()

	log.Println("requests ended")
}
