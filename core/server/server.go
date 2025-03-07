package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/core/base"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	addr string

	e          *gin.Engine
	middleware []gin.HandlerFunc
	routers    []func(*gin.RouterGroup)
}

func New(cnf *Config) *Server {
	if !cnf.debug {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		addr:       fmt.Sprintf("%s:%d", cnf.host, cnf.port),
		e:          gin.New(),
		middleware: make([]gin.HandlerFunc, 0),
		routers:    make([]func(*gin.RouterGroup), 0),
	}

	// 注册基础路由
	s.Routers(base.Router)

	return s
}

func (s *Server) Middleware(fn ...gin.HandlerFunc) *Server {
	s.middleware = append(s.middleware, fn...)
	return s
}

func (s *Server) Routers(fn ...func(*gin.RouterGroup)) *Server {
	s.routers = append(s.routers, fn...)
	return s
}

func (s *Server) GC() {
	s.routers = nil
}

func (s *Server) Run() error {
	// 注册中间件
	s.e.Use(s.middleware...)
	// 注册路由
	for _, fn := range s.routers {
		fn(s.e.Group("/"))
	}
	s.GC()
	return nil
}

func (s *Server) ListenAndServer() {
	if err := s.Run(); err != nil {
		log.Fatalln(err)
	}
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
