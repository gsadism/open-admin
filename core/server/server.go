package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	addr string

	e *gin.Engine
}

func New(cnf *Config) *Server {
	if !cnf.debug {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		addr: fmt.Sprintf("%s:%d", cnf.host, cnf.port),
		e:    gin.New(),
	}

	return s
}

func (s *Server) GC() {

}

func (s *Server) ListenAndServer() {
	srv := http.Server{
		Addr:    s.addr,
		Handler: s.e,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("%s\n", err)
		}
	}()
	log.Printf("Listen: %s\n", s.addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
