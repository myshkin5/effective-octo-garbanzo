package handlers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
)

const (
	ok         = "{}"
	notAllowed = `{
		"code": 405,
		"error": "Method not allowed",
		"status": "Method Not Allowed"
	}`
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(ok))
})

var _ = Describe("Health", func() {
	var (
		recorder *httptest.ResponseRecorder
		request  *http.Request
		router   *mux.Router
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		router = mux.NewRouter()
		handlers.MapHealthRoutes(router, alice.Chain{})
	})

	Describe("happy path", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest(http.MethodGet, "/health", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a good health body", func() {
			tests := []struct {
				req     *http.Request
				handler http.Handler
				code    int
				allow   string // Contents of the Allow header
				body    string
			}{
				// No handlers
				{newRequest("GET", "/foo"), handlers.MethodHandler{}, http.StatusMethodNotAllowed, "", notAllowed},
				{newRequest("OPTIONS", "/foo"), handlers.MethodHandler{}, http.StatusOK, "", ""},

				// A single handler
				{newRequest("GET", "/foo"), handlers.MethodHandler{"GET": okHandler}, http.StatusOK, "", ok},
				{newRequest("POST", "/foo"), handlers.MethodHandler{"GET": okHandler}, http.StatusMethodNotAllowed, "GET", notAllowed},

				// Multiple handlers
				{newRequest("GET", "/foo"), handlers.MethodHandler{"GET": okHandler, "POST": okHandler}, http.StatusOK, "", ok},
				{newRequest("POST", "/foo"), handlers.MethodHandler{"GET": okHandler, "POST": okHandler}, http.StatusOK, "", ok},
				{newRequest("DELETE", "/foo"), handlers.MethodHandler{"GET": okHandler, "POST": okHandler}, http.StatusMethodNotAllowed, "GET, POST", notAllowed},
				{newRequest("OPTIONS", "/foo"), handlers.MethodHandler{"GET": okHandler, "POST": okHandler}, http.StatusOK, "GET, POST", ""},

				// Override OPTIONS
				{newRequest("OPTIONS", "/foo"), handlers.MethodHandler{"OPTIONS": okHandler}, http.StatusOK, "", ok},
			}

			for i, test := range tests {
				rec := httptest.NewRecorder()
				test.handler.ServeHTTP(rec, test.req)
				Expect(rec.Code).To(Equal(test.code), "%d: wrong code, got %d want %d", i, rec.Code, test.code)
				allow := rec.HeaderMap.Get("Allow")
				Expect(allow).To(Equal(test.allow), "%d: wrong Allow, got %s want %s", i, allow, test.allow)
				body := rec.Body.String()
				if len(body) == 0 {
					Expect(body).To(Equal(test.body), "%d: wrong body, got %q want %q", i, body, test.body)
				} else {
					Expect(body).To(MatchJSON(test.body), "%d: wrong body, got %q want %q", i, body, test.body)
				}
			}
		})
	})
})

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}
