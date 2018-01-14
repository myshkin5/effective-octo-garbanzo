package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/garbanzo"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/octo"
)

const (
	samples = 3
	count   = 100
	url     = "http://localhost:8080/"
)

var _ = Describe("API", func() {
	var (
		database persistence.Database
		token    string
	)

	BeforeSuite(func() {
		var err error
		database, err = persistence.Open()
		Expect(err).NotTo(HaveOccurred())

		_, err = database.Exec(context.Background(), "insert into org (name) values ('org1') returning id")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		_, err := database.Exec(context.Background(), "delete from garbanzo")
		Expect(err).NotTo(HaveOccurred())
		_, err = database.Exec(context.Background(), "delete from octo")
		Expect(err).NotTo(HaveOccurred())
		_, err = database.Exec(context.Background(), "delete from org where name = 'org1'")
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		response, err := http.Get("http://localhost:8081/token")
		Expect(err).NotTo(HaveOccurred())

		defer response.Body.Close()

		Expect(response.StatusCode).To(Equal(http.StatusOK))

		bytes, err := ioutil.ReadAll(response.Body)
		Expect(err).NotTo(HaveOccurred())

		var body handlers.JSONObject
		err = json.Unmarshal(bytes, &body)
		Expect(err).NotTo(HaveOccurred())

		token = "bearer " + body["token"].(string)
	})

	Measure("the standard suite of operations", func(b Benchmarker) {
		b.Time("runtime", func() {
			errs := make(chan error, samples*count*2)

			var wg sync.WaitGroup
			wg.Add(6)
			go postOctos(token, errs, &wg)
			go deleteOctos(token, errs, &wg)
			go getOctos(token, errs, &wg)
			go postGarbanzos(token, errs, &wg)
			go deleteGarbanzos(token, errs, &wg)
			go getGarbanzos(token, errs, &wg)

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

func postOctos(token string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	nanos := time.Now().UnixNano()

	for i := 0; i < count; i++ {
		body := fmt.Sprintf(`{
			"name": "stress_%d_%d"
		}`, nanos, i)
		response, err := do("POST", url+"octos", token, strings.NewReader(body))
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

func deleteOctos(token string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		octos := getRandomOctos(10, url, token, errs)

		for _, octo := range octos {
			response, err := do("DELETE", octo.Link, token, nil)
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

func getOctos(token string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		octos := getRandomOctos(10, url, token, errs)

		for _, octo := range octos {
			response, err := do("GET", octo.Link, token, nil)
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

func postGarbanzos(token string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		octos := getRandomOctos(10, url, token, errs)

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
			response, err := do("POST", octo.Garbanzos, token, strings.NewReader(body))
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

func deleteGarbanzos(token string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		garbanzos := getRandomGarbanzos(10, url, token, errs)

		for _, garbanzo := range garbanzos {
			response, err := do("DELETE", garbanzo.Link, token, nil)
			if err != nil {
				errs <- err
				continue
			}
			if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusNotFound {
				errs <- fmt.Errorf("deleting garbanzos expecting status %d, got %d", http.StatusNoContent, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}
}

func getGarbanzos(token string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < count/10; i++ {
		garbanzos := getRandomGarbanzos(10, url, token, errs)

		for _, garbanzo := range garbanzos {
			response, err := do("GET", garbanzo.Link, token, nil)
			if err != nil {
				errs <- err
				continue
			}
			if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
				errs <- fmt.Errorf("getting garbanzos expecting status %d, got %d", http.StatusOK, response.StatusCode)
				continue
			}

			err = response.Body.Close()
			if err != nil {
				errs <- err
			}
		}
	}
}

func getRandomOctos(count int, url, token string, errs chan error) []octo.Octo {
	response, err := do("GET", url+"octos", token, nil)
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

func getRandomGarbanzos(count int, url, token string, errs chan error) []garbanzo.Garbanzo {
	octoList := getRandomOctos(count, url, token, errs)

	var list []garbanzo.Garbanzo
	for _, octo := range octoList {
		response, err := do("GET", octo.Garbanzos, token, nil)
		if err != nil {
			errs <- err
			return nil
		}

		if response.StatusCode != http.StatusOK {
			errs <- fmt.Errorf("getting garbanzos expecting status %d, got %d", http.StatusOK, response.StatusCode)
			return nil
		}

		var fullList []garbanzo.Garbanzo
		err = json.NewDecoder(response.Body).Decode(&fullList)
		if err != nil {
			errs <- err
			return nil
		}

		err = response.Body.Close()
		if err != nil {
			errs <- err
		}

		if len(fullList) > 0 {
			list = append(list, fullList[rand.Intn(len(fullList))])
		}
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

func do(method, url, token string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", token)
	return http.DefaultClient.Do(request)
}
