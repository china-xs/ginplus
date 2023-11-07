// Package http
// @file      : http.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/7 17:30
// @Description:
package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net"
	"net/http"
)

const OperationKey = "operation"

// ServerOption is an HTTP server option.
type ServerOption func(*Server)

type Server struct {
	*gin.Engine
	isDebug bool               // 测试环境
	addr    string             // default 0.0.0.0:8080
	dec     DecodeRequestFunc  // 请求参数绑定结构
	enc     EncodeResponseFunc // 定义返回结构
	ms      []Middleware       // 全局中间价
	filter  []gin.HandlerFunc  // gin 全局中间件， 执行比ms 早
	ops     []ServerOption
}

// Start http server start
func (this *Server) Start(ctx context.Context) (err error) {
	r := gin.New()
	srv := &Server{
		Engine: r,
		addr:   "0.0.0.0:8080",
		dec:    DefaultRequestDecoder,
		enc:    DefaultResponseEncoder,
	}
	for _, o := range this.ops {
		o(srv)
	}
	// use gin middleware
	if len(srv.filter) > 0 {
		srv.Engine.Use(srv.filter...)
	}
	s := http.Server{
		Addr:    this.addr,
		Handler: this,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			ctx = context.WithValue(ctx, `debug`, `sss`)
			return ctx
		},
	}
	err = s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	// 暂时不加tls
	return
}

// Stop http server stop
func (this Server) Stop(ctx context.Context) (err error) {
	return
}

// Middleware 引用中间件
func (this Server) Middleware(h Handler) Handler {
	return Chain(this.ms...)(h)
}

// Bind  请求参数绑定
func (s *Server) Bind(c *gin.Context, obj any) error {
	return s.dec(c, obj)
}

// Result 结果结返回
func (s *Server) Result(c *gin.Context, obj any, err error) {
	s.enc(c, obj, err)
	return
}
