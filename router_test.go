package uRouter

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	var as = assert.New(t)

	t.Run("", func(t *testing.T) {
		var list []int
		var r = New()
		r.Use(func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})

		var g0 = r.Group("", func(ctx *Context) {
			list = append(list, 3)
			ctx.Next()
			list = append(list, 4)
		})

		var g1 = g0.Group("/api/v1", func(ctx *Context) {
			list = append(list, 5)
			ctx.Next()
			list = append(list, 6)
		})

		g1.On("greet", func(ctx *Context) {
			list = append(list, 9)
		}, func(ctx *Context) {
			list = append(list, 7)
			ctx.Next()
			list = append(list, 8)
		})

		r.Display()

		r.Emit(&Context{
			Request: &Request{
				Header: NewHttpHeader(http.Header{"X-Path": []string{"/api/v1/greet"}}),
				Body:   nil,
			},
			Writer: nil,
		})

		as.Equal(9, len(list))
		as.Equal(1, list[0])
		as.Equal(3, list[1])
		as.Equal(5, list[2])
		as.Equal(7, list[3])
		as.Equal(9, list[4])
		as.Equal(8, list[5])
		as.Equal(6, list[6])
		as.Equal(4, list[7])
		as.Equal(2, list[8])
	})

	t.Run("", func(t *testing.T) {
		var r = New()
		var list []int

		r.On("test", func(ctx *Context) {
			list = append(list, 3)
		}, func(ctx *Context) {
			list = append(list, 1)
			ctx.Next()
			list = append(list, 2)
		})

		r.Emit(&Context{
			Request: &Request{
				Header: NewHttpHeader(http.Header{"X-Path": []string{"/test"}}), Body: nil,
			},
			Writer: nil,
		})

		as.Equal(len(list), 3)
		as.Equal(1, list[0])
		as.Equal(3, list[1])
		as.Equal(2, list[2])
	})

	t.Run("no match", func(t *testing.T) {
		var r = New()
		var list []int

		r.OnNoMatch = func(ctx *Context) {
			list = append(list, 1)
		}

		r.Emit(&Context{
			Request: &Request{
				Header: NewHttpHeader(http.Header{XPath: []string{"/test"}}), Body: nil,
			},
			Writer: nil,
		})

		as.Equal(len(list), 1)
		as.Equal(1, list[0])
	})

	t.Run("no handler", func(t *testing.T) {
		var r = New()
		var list []int

		r.Group("test")

		r.Emit(&Context{
			Request: &Request{
				Header: NewHttpHeader(http.Header{XPath: []string{"/test"}}), Body: nil,
			},
			Writer: nil,
		})

		r.routes["404"] = nil

		r.Display()
		as.Equal(len(list), 0)
	})
}