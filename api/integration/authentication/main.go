package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

func main() {
	initLogging()

	router := initRoutes()

	serverAddr := persistence.GetEnvWithDefault("SERVER_ADDR", "localhost")
	port := persistence.GetEnvWithDefault("PORT", "8081")
	listenAndServe(serverAddr, port, router)
}

func initLogging() {
	err := logs.Init()
	if err != nil {
		panic(err)
	}
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()

	middleware := alice.New(handlers.LoggingHandler)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logs.Logger.Panic("Could not generate key: ", err)
	}

	eBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(eBuf, uint32(privateKey.PublicKey.E))

	router.PathPrefix("/keys").Handler(middleware.ThenFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlers.Respond(w, http.StatusOK, handlers.JSONObject{
			"keys": []handlers.JSONObject{
				{
					"alg": "RS256",
					"kid": "the-one-and-only",
					"kty": "RSA",
					"n":   base64.URLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
					"e":   base64.URLEncoding.EncodeToString(eBuf)[:4],
					"use": "sig",
				},
			},
		})
	}))

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 300,
	}
	token := &jwt.Token{
		Header: map[string]interface{}{
			"kid": "the-one-and-only",
			"alg": jwt.SigningMethodRS256.Alg(),
		},
		Claims: claims,
		Method: jwt.SigningMethodRS256,
	}
	signedString, err := token.SignedString(privateKey)
	if err != nil {
		logs.Logger.Panic("Could not sign token: ", err)
	}

	router.PathPrefix("/token").Handler(middleware.ThenFunc(func(w http.ResponseWriter, _ *http.Request) {
		handlers.Respond(w, http.StatusOK, handlers.JSONObject{
			"token": signedString,
		})
	}))

	return router
}

func listenAndServe(serverAddr, port string, router *mux.Router) {
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	err := http.ListenAndServe(serverAddr+":"+port, router)
	if err != nil {
		logs.Logger.Panic("ListenAndServe: ", err)
	}
}
