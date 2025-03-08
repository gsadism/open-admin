package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/core/model"
	"github.com/gsadism/open-admin/pkg/file"
	"github.com/panjf2000/ants/v2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type StaticFile struct {
	RelativePath string
	FilePath     string
}

type Server struct {
	addr string

	pSize int

	e          *gin.Engine
	middleware []gin.HandlerFunc
	routers    []func(*gin.RouterGroup)
	files      []StaticFile

	hooks []func()
}

func New(cnf *Config) *Server {
	if !cnf.debug {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		addr:       fmt.Sprintf("%s:%d", cnf.host, cnf.port),
		pSize:      2,
		e:          gin.New(),
		middleware: make([]gin.HandlerFunc, 0),
		routers:    make([]func(*gin.RouterGroup), 0),
		files:      make([]StaticFile, 0),
		hooks:      make([]func(), 0),
	}

	// 注册基础路由
	// s.Routers(base.Router)

	return s
}

func (s *Server) SetPoolSize(size int) *Server {
	s.pSize = size
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

func (s *Server) Files(RelativePath, FilePath string) *Server {
	if !filepath.IsAbs(FilePath) {
		if d, err := filepath.Abs(FilePath); err != nil {
			log.Println(err.Error())
		} else {
			FilePath = d
		}
	}
	if file.Exists(FilePath) {
		s.files = append(s.files, StaticFile{
			RelativePath: RelativePath,
			FilePath:     FilePath,
		})
	}
	return s
}

func (s *Server) Hoos(fn ...func()) *Server {
	s.hooks = append(s.hooks, fn...)
	return s
}

func (s *Server) GC() {
	s.routers = nil
	s.middleware = nil
	s.files = nil
	s.hooks = nil
}

func (s *Server) AutoMigrate(client *gorm.DB, auto bool, models ...model.IModel) *Server {
	if auto {
		for _, table := range models {
			_ = client.AutoMigrate(table)
		}
	}

	return s
}

func (s *Server) RunHooks() {
	if p, err := ants.NewPool(s.pSize, ants.WithPreAlloc(true)); err != nil {
		log.Fatal(err.Error())
	} else {
		defer p.Release()

		var wg sync.WaitGroup
		for _, fn := range s.hooks {
			wg.Add(1)
			_ = ants.Submit(func() {
				defer wg.Done()
				fn()
			})
		}
		wg.Wait()
	}
}

func (s *Server) Run() error {
	// 执行预加载hooks
	s.RunHooks()
	// 注册中间件
	s.e.Use(s.middleware...)
	// 注册static file
	for _, f := range s.files {
		s.e.StaticFile(f.RelativePath, f.FilePath)
	}
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
