package uRouter

import (
	"net/http"
	"testing"
)

func BenchmarkOneRoute(b *testing.B) {
	router := New()
	router.OnGET("/ping", func(c *Context) {
	})
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/ping", ctx)
	}
}

func BenchmarkOneRouteDynamic(b *testing.B) {
	router := New()
	router.OnGET("/user/:id", func(c *Context) {
	})
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/user/1", ctx)
	}
}

func BenchmarkRecoveryMiddleware(b *testing.B) {
	router := New()
	router.Use(Recovery())
	router.OnGET("/", func(c *Context) {
	})
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/", ctx)
	}
}

func Benchmark5Params(b *testing.B) {
	router := New()
	router.Use(func(ctx *Context) {})
	router.OnGET("/param/:param1/:params2/:param3/:param4/:param5", func(c *Context) {
	})
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/param/path/to/parameter/john/12345", ctx)
	}
}

func BenchmarkOneRouteJSON(b *testing.B) {
	router := New()
	router.Use(func(ctx *Context) {})
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.OnGET("/json", func(c *Context) {
		_ = c.WriteJSON(http.StatusOK, data)
		c.Request.Close()
	})
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/json", ctx)
	}
}

func Benchmark404(b *testing.B) {
	router := New()
	router.OnGET("/", func(c *Context) {})
	router.OnGET("/something", func(c *Context) {})
	router.OnNotFound = func(ctx *Context) {}
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/ping", ctx)
	}
}

func Benchmark404Many(b *testing.B) {
	router := New()
	router.OnGET("/", func(c *Context) {})
	router.OnGET("/path/to/something", func(c *Context) {})
	router.OnGET("/post/:id", func(c *Context) {})
	router.OnGET("/view/:id", func(c *Context) {})
	router.OnGET("/favicon.ico", func(c *Context) {})
	router.OnGET("/robots.txt", func(c *Context) {})
	router.OnGET("/delete/:id", func(c *Context) {})
	router.OnGET("/user/:id/:mode", func(c *Context) {})
	router.OnNotFound = func(ctx *Context) {}
	router.StartSilently()
	ctx := newContextMocker()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.EmitEvent("GET", "/viewfake", ctx)
	}
}
