package integration_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/octo"
)

var _ = Describe("API", func() {
	quickCreateOctos := func(count int, url string, errs chan error, wg *sync.WaitGroup) {
		defer wg.Done()

		nanos := time.Now().UnixNano()

		for i := 0; i < count; i++ {
			body := fmt.Sprintf(`{
					"name": "stress_%d_%d"
				}`, nanos, i)
			response, err := http.Post(url+"octos", "application/json", strings.NewReader(body))
			if err != nil {
				errs <- err
				continue
			}

			if response.StatusCode != http.StatusCreated {
				errs <- fmt.Errorf("create octos expecting status %d, got %d", http.StatusCreated, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}

	quickDeleteRandomOctos := func(count int, url string, errs chan error, wg *sync.WaitGroup) {
		defer wg.Done()

		for i := 0; i < count; i++ {
			response, err := http.Get(url + "octos")
			if err != nil {
				errs <- err
				continue
			}

			if response.StatusCode != http.StatusOK {
				errs <- fmt.Errorf("getting octos expecting status %d, got %d", http.StatusOK, response.StatusCode)
				continue
			}

			var list []octo.Octo
			err = json.NewDecoder(response.Body).Decode(&list)
			if err != nil {
				errs <- err
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}

			if len(list) == 0 {
				continue
			}

			octo := list[rand.Intn(len(list))]

			request, err := http.NewRequest("DELETE", octo.Link, nil)
			if err != nil {
				errs <- err
				continue
			}
			response, err = http.DefaultClient.Do(request)
			if err != nil {
				errs <- err
				continue
			}
			if response.StatusCode != http.StatusNoContent {
				errs <- fmt.Errorf("deleting octos expecting status %d, got %d", http.StatusNoContent, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}

	Measure("the standard suite of operations", func(b Benchmarker) {
		b.Time("runtime", func() {
			url := "http://localhost:8080/"
			errs := make(chan error, 2000)

			var wg sync.WaitGroup
			wg.Add(2)
			go quickCreateOctos(100, url, errs, &wg)
			go quickDeleteRandomOctos(100, url, errs, &wg)

			wg.Wait()
			close(errs)

			count := len(errs)

			for err := range errs {
				fmt.Printf("Test received error: %v\n", err)
			}

			Expect(count).To(BeZero())
		})
	}, 3)
})
