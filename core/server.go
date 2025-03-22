package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/logging"
	"github.com/gsadism/open-admin/pkg/object"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Folder : 获取路径
func Folder(path string) string {
	if !filepath.IsAbs(path) {
		if d, err := filepath.Abs(path); err != nil {
			Exit(err.Error())
		} else {
			path = d
		}
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// 创建路径
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				Exit(err.Error())
			}
		}
	}
	return path
}

type Server struct {
	host string
	port int

	engine     *gin.Engine
	middleware []gin.HandlerFunc
	routers    []func(*gin.RouterGroup)
	files      map[string]string
}

func New(v *viper.Viper) *Server {
	// 日志记录器初始化
	logging.ReplaceGlobals(logging.New().SetSkip(2).File(
		Folder(v.GetString("logger.file.path")),
		object.Default[string](v.GetString("logger.file.name"), "open-admin.log"),
		v.GetString("logger.file.level"),
		v.GetInt("logger.file.max-size"),
		v.GetInt("logger.file.max-age"),
		v.GetInt("logger.file.max-backups"),
		v.GetBool("logger.file.compress"),
	).R())

	if !v.GetBool("server.debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		host: parseIP(v.GetString("server.host"), defaultHost),
		port: parsePort(v.GetInt("server.port"), defaultPort),

		engine:     gin.New(),
		middleware: make([]gin.HandlerFunc, 0),
		routers:    make([]func(*gin.RouterGroup), 0),
		files:      make(map[string]string),
	}

	s.Routers(routers).
		StaticFile("favicon.ico", v.GetString("website.favicon-ico"))

	return s
}

func Default() *Server {
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		host: defaultHost,
		port: defaultPort,

		engine:     gin.New(),
		middleware: make([]gin.HandlerFunc, 0),
		routers:    make([]func(*gin.RouterGroup), 0),
		files:      make(map[string]string),
	}
	s.Routers(routers)

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

func (s *Server) StaticFile(relative, path string) *Server {
	//s.engine.StaticFile()
	if !filepath.IsAbs(path) {
		if d, err := filepath.Abs(path); err != nil {
			return s
		} else {
			if _, err := os.Stat(d); err != nil {
				if os.IsNotExist(err) {
					return s
				}
				return s
			} else {
				s.files[relative] = d
			}
		}
	}
	return s
}

func (s *Server) run() error {

	// 注册全局中间件
	s.engine.Use(s.middleware...)

	// 注入 static files
	for relative, path := range s.files {
		s.engine.StaticFile(relative, path)
	}

	// 注入router
	for _, fn := range s.routers {
		fn(s.engine.Group("/"))
	}

	return nil
}

func (s *Server) gc() {
	s.middleware = nil
	s.routers = nil
	s.files = nil
}

func (s *Server) ListenAndServer() {
	if err := s.run(); err != nil {
		Exit(err.Error())
	} else {
		srv := http.Server{
			Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
			Handler: s.engine,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
				Exit(err.Error())
			}
		}()
		Debug(fmt.Sprintf("Listen: %s\n", fmt.Sprintf("%s:%d", s.host, s.port)))
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			Exit(fmt.Sprint("Server Shutdown:", err))
		}
	}
}
