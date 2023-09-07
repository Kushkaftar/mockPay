package app

import "net/http"

type server struct {
	httpServer *http.Server
}

func (s *server) serverRun(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}
