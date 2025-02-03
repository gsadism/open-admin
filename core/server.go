package core

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gsadism/open-admin/depends/logging"
	"github.com/gsadism/open-admin/osv"
	"github.com/gsadism/open-admin/pkg/crypto/_rsa"
	"github.com/gsadism/open-admin/pkg/file"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Options struct {
	Debug bool

	Host string
	Port int

	Resource string
}

type Server struct {
	addr string

	gin        *gin.Engine
	routers    []func(*gin.RouterGroup)
	middleware []gin.HandlerFunc

	resource string
}

func NewServer(opt *Options) *Server {
	if !opt.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &Server{
		gin:        gin.New(),
		routers:    make([]func(*gin.RouterGroup), 0),
		middleware: make([]gin.HandlerFunc, 0),

		resource: opt.Resource,
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

func (s *Server) writeIcon() {
	b64Img := icon
	if idx := strings.Index(b64Img, ","); idx != -1 {
		b64Img = b64Img[idx+1:]
	}
	// 解码
	if img, err := base64.StdEncoding.DecodeString(b64Img); err != nil {
		logging.Fatal(err.Error())
	} else {
		if err := ioutil.WriteFile(filepath.Join(s.resource, "favicon.ico"), img, 0644); err != nil {
			logging.Fatal(err.Error())
		}
	}
}

func (s *Server) writeRsaFile() (*rsa.PublicKey, *rsa.PrivateKey) {
	if !file.Exists(filepath.Join(s.resource, "public_key.pem")) || !file.Exists(filepath.Join(s.resource, "private_key.pem")) {
		if file.Exists(filepath.Join(s.resource, "public_key.pem")) {
			if err := os.Remove(filepath.Join(s.resource, "public_key.pem")); err != nil {
				logging.Fatal(err.Error())
			}
		}
		if file.Exists(filepath.Join(s.resource, "private_key.pem")) {
			if err := os.Remove(filepath.Join(s.resource, "private_key.pem")); err != nil {
				logging.Fatal(err.Error())
			}
		}
		prv, pub := _rsa.Generate(_rsa.MEDIUM)
		// 写入pem
		if err := _rsa.WritePKIXPublicKeyFile(pub, filepath.Join(s.resource, "public_key.pem"), nil); err != nil {
			logging.Fatal(err.Error())
		}
		if err := _rsa.WritePKCS8PrivateKeyFile(prv, filepath.Join(s.resource, "private_key.pem"), nil); err != nil {
			logging.Fatal(err.Error())
		}
		return pub, prv
	} else {
		var prv *rsa.PrivateKey
		var pub *rsa.PublicKey
		if d, err := _rsa.LoadPKIXPublicKeyWithFile(filepath.Join(s.resource, "public_key.pem")); err != nil {
			logging.Fatal(err.Error())
		} else {
			pub = d
		}
		if d, err := _rsa.LoadPKCS8PrivateKeyWithFile(filepath.Join(s.resource, "private_key.pem")); err != nil {
			logging.Fatal(err.Error())
		} else {
			prv = d
		}
		return pub, prv
	}
}

func (s *Server) run() {
	//TODO 注册全局中间件
	for _, m := range s.middleware {
		s.gin.Use(m)
	}

	if !file.Exists(s.resource) {
		if err := os.MkdirAll(s.resource, os.ModePerm); err != nil {
			logging.Fatal(err.Error())
		}
	}
	// rsa
	osv.PublicKey, osv.PrivateKey = s.writeRsaFile()
	// 注册 favicon.ico
	if !file.Exists(filepath.Join(s.resource, "favicon.ico")) {
		s.writeIcon()
	}
	s.gin.StaticFile("favicon.ico", filepath.Join(s.resource, "favicon.ico"))
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
			//fmt.Println(fmt.Sprintf("\033[%dm[%v] %v\033[0m", 30+1, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
			//os.Exit(-1)
			logging.Fatal(err.Error())
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
		//fmt.Println(fmt.Sprintf("\033[%dm[%v] %v\033[0m", 30+1, time.Now().Format("2006-01-02 15:04:05"), err.Error()))
		//os.Exit(-1)
		logging.Fatal(err.Error())
	}
}
