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

const (
	samples = 3
	count   = 100
)

var _ = Describe("API", func() {
	Measure("the standard suite of operations", func(b Benchmarker) {
		b.Time("runtime", func() {
			url := "http://localhost:8080/"
			errs := make(chan error, samples*count*2)

			var wg sync.WaitGroup
			wg.Add(4)
			go quickCreateOctos(count, url, errs, &wg)
			go quickDeleteRandomOctos(count, url, errs, &wg)
			go quickGetRandomOctos(count, url, errs, &wg)
			go quickCreateRandomGarbanzos(count, url, errs, &wg)

			timedOut := waitTimeout(&wg, time.Minute)
			Expect(timedOut).To(BeFalse())
			close(errs)

			errCount := len(errs)

			for err := range errs {
				fmt.Printf("Test received error: %v\n", err)
			}

			Expect(errCount).To(BeZero())
		})
	}, samples)
})

func quickCreateOctos(count int, url string, errs chan error, wg *sync.WaitGroup) {
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

func quickDeleteRandomOctos(count int, url string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		octos := getRandomOctos(10, url, errs)

		for _, octo := range octos {
			request, err := http.NewRequest("DELETE", octo.Link, nil)
			if err != nil {
				errs <- err
				continue
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				errs <- err
				continue
			}
			if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusNotFound {
				errs <- fmt.Errorf("deleting octos expecting status %d, got %d", http.StatusNoContent, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}
}

func quickGetRandomOctos(count int, url string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		octos := getRandomOctos(10, url, errs)

		for _, octo := range octos {
			response, err := http.Get(octo.Link)
			if err != nil {
				errs <- err
				continue
			}
			if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
				errs <- fmt.Errorf("getting octos expecting status %d, got %d", http.StatusOK, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}
}

func quickCreateRandomGarbanzos(count int, url string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		octos := getRandomOctos(10, url, errs)

		for j, octo := range octos {
			var garbanzoType string
			if j%2 == 0 {
				garbanzoType = "DESI"
			} else {
				garbanzoType = "KABULI"
			}
			body := fmt.Sprintf(`{
				"type": "%s",
				"diameter-mm": 2.0
			}`, garbanzoType)
			response, err := http.Post(octo.Garbanzos, "application/json", strings.NewReader(body))
			if err != nil {
				errs <- err
				continue
			}

			if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusConflict {
				errs <- fmt.Errorf("create garbanzos expecting status %d, got %d", http.StatusCreated, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}
}

func getRandomOctos(count int, url string, errs chan error) []octo.Octo {
	response, err := http.Get(url + "octos")
	if err != nil {
		errs <- err
		return nil
	}

	if response.StatusCode != http.StatusOK {
		errs <- fmt.Errorf("getting octos expecting status %d, got %d", http.StatusOK, response.StatusCode)
		return nil
	}

	var fullList []octo.Octo
	err = json.NewDecoder(response.Body).Decode(&fullList)
	if err != nil {
		errs <- err
		return nil
	}

	err = response.Body.Close()
	if err != nil {
		errs <- err
	}

	returnCount := min(count, len(fullList))

	list := make([]octo.Octo, returnCount)
	for i := 0; i < returnCount; i++ {
		// There may be duplicates
		list[i] = fullList[rand.Intn(len(fullList))]
	}

	return list
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
