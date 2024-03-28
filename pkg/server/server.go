package server

import (
	"fmt"
	"log"
	"net/http"
	"qrdb/pkg/di"
	"qrdb/pkg/server/routes"
	"qrdb/pkg/server/routes/qrcode"
)

type server struct {
	dependencies         di.Dependencies
	pathToRoutable       map[string]Routable
	pathToAllowedMethods map[string][]string
	pathToAllowedRoles   map[string][]string
}

func NewServer() *server {
	return &server{
		dependencies:         di.GetDependencies(),
		pathToRoutable:       make(map[string]Routable),
		pathToAllowedMethods: make(map[string][]string),
		pathToAllowedRoles:   make(map[string][]string),
	}
}

func (s *server) Start() {

	OPEN := []string{}
	ALL_METHODS := []string{}

	s.RegisterRoute(
		"/status",
		routes.Status{},
		ALL_METHODS,
		OPEN,
	)

	s.RegisterRoute(
		"/qr",
		qrcode.CreateQRCode{},
		[]string{"POST"},
		OPEN,
	)

	s.RegisterRoute(
		"/qr/:code",
		qrcode.GetData{},
		[]string{"GET"},
		OPEN,
	)

	s.RegisterRoute(
		"/qr/render/:code",
		qrcode.RenderQR{},
		[]string{"GET"},
		OPEN,
	)

	s.RegisterRoute(
		"/qr/render/:code/page/:page/of/:max_page",
		qrcode.RenderQR{},
		[]string{"GET"},
		OPEN,
	)

	log.Print("Listening...")
	if err := http.ListenAndServe(":8080", s); err != nil {
		fmt.Println(err)
	}
}
