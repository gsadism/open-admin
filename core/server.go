package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	Debug bool

	Host string
	Port int
}

type Server struct {
	addr string

	gin        *gin.Engine
	routers    []func(*gin.RouterGroup)
	middleware []gin.HandlerFunc
}

func NewServer(opt *Options) *Server {
	if !opt.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		gin:        gin.New(),
		routers:    make([]func(*gin.RouterGroup), 0),
		middleware: make([]gin.HandlerFunc, 0),
	}
	s.setAddr(opt.Host, opt.Port)

	return s
}

func (s *Server) Routers(r ...func(*gin.RouterGroup)) *Server {
	s.routers = append(s.routers, r...)
	return s
}

func (s *Server) Middleware(m ...gin.HandlerFunc) *Server {
	s.middleware = append(s.middleware, m...)
	return s
}

func (s *Server) setAddr(host string, port int) {
	s.addr = fmt.Sprintf("%s:%d", func(ip string) string {
		if net.ParseIP(ip) == nil {
			return "0.0.0.0"
		}
		return ip
	}(host), func(port int) int {
		if port < 0 || port > 65535 {
			return 9815
		}
		return port
	}(port))
}

func (s *Server) run() {
	//TODO 注册全局中间件
	for _, m := range s.middleware {
		s.gin.Use(m)
	}

	//TODO 注册路由
	for _, router := range s.routers {
		router(s.gin.Group(""))
	}
}

func (s *Server) clear() {
	s.routers = nil
	s.middleware = nil
}

func (s *Server) ListenAndServe() {
	s.run()
	s.clear()

	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.gin,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Println(fmt.Sprintf("\033[%dm[%v] %v\033[0m", 30+1, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
			os.Exit(-1)
		}
	}()
	fmt.Println(fmt.Sprintf("\033[%dm[%v] %v\033[0m", 30+4, time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Server Listen %v", srv.Addr)))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println(fmt.Sprintf("\033[%dm[%v] %v\033[0m", 30+1, time.Now().Format("2006-01-02 15:04:05"), "Server Shutdown..."))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println(fmt.Sprintf("\033[%dm[%v] %v\033[0m", 30+1, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		os.Exit(-1)
	}
}
