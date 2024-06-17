package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type WebServer struct {
	Router       chi.Router
	Routes       []Route
	ServiceName  string
	WebServerUrl string
}

func NewWebServer(webServerUrl, serviceName string, routes []Route) *WebServer {
	webServer := &WebServer{
		Router:       chi.NewRouter(),
		Routes:       routes,
		ServiceName:  serviceName,
		WebServerUrl: webServerUrl,
	}

	webServer.Router.Use(middleware.RequestID)
	webServer.Router.Use(middleware.RealIP)
	webServer.Router.Use(middleware.Recoverer)
	webServer.Router.Use(middleware.Logger)
	webServer.Router.Handle("/metrics", promhttp.Handler())

	for _, route := range webServer.Routes {
		webServer.Router.Method(route.Method, route.Path, route.Handler)
	}

	return webServer
}

func (s *WebServer) Start() {
	fmt.Printf("Starting %s webserver on %s\n", s.ServiceName, s.WebServerUrl)
	log.Printf("Starting %s webserver on %s\n", s.ServiceName, s.WebServerUrl)
	err := http.ListenAndServe(s.WebServerUrl, s.Router)
	if err != nil {
		log.Fatal(err)
	}
}
