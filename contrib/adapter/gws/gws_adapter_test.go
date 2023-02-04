package gws

import (
	"bytes"
	"github.com/lxzan/gws"
	"github.com/lxzan/uRouter"
	"github.com/stretchr/testify/assert"
	"testing"
)

type connMocker struct {
	opcode gws.Opcode
	buf    *bytes.Buffer
}

func (c *connMocker) WriteMessage(opcode gws.Opcode, payload []byte) {
	c.opcode = opcode
	c.buf.Write(payload)
}

func TestNewAdapter(t *testing.T) {
	var as = assert.New(t)

	t.Run("normal", func(t *testing.T) {
		const requestPayload = "hello"
		const responsePayload = "world"
		var sum = 0
		var router = uRouter.New()
		var adapter = NewAdapter(router, uRouter.TextHeader)

		router.On("testEncode", func(ctx *uRouter.Context) {
			ctx.Writer = &responseWriter{
				conn:        &connMocker{buf: bytes.NewBufferString("")},
				headerCodec: uRouter.TextHeader,
				header:      uRouter.TextHeader.Generate(),
			}

			sum++
			ctx.Writer.Header().Set(uRouter.ContentType, uRouter.MimeStream)
			ctx.Writer.Header().Set(uRouter.XPath, "/testDecode")
			ctx.Writer.Code(int(gws.OpcodeText))
			_, _ = ctx.Writer.Write([]byte(responsePayload))
			ctx.Writer.RawResponseWriter()
			if err := ctx.Writer.Flush(); err != nil {
				as.NoError(err)
				return
			}

			as.Equal(2, ctx.Request.Header.Len())
			as.Equal(requestPayload, ctx.Request.Body.(*bytes.Buffer).String())

			var writer = ctx.Writer.RawResponseWriter().(*connMocker)
			if err := adapter.ServeWebSocket(nil, writer.buf.Bytes()); err != nil {
				as.NoError(err)
				return
			}
		})

		router.On("testDecode", func(ctx *uRouter.Context) {
			sum++
			as.Equal(2, ctx.Request.Header.Len())
			as.Equal(responsePayload, ctx.Request.Body.(*bytes.Buffer).String())
		})

		var b = bytes.NewBufferString("")
		var header = uRouter.F{
			uRouter.ContentType: uRouter.MimeJson,
			uRouter.XPath:       "/testEncode",
		}
		if err := adapter.codec.Encode(b, header); err != nil {
			as.NoError(err)
			return
		}
		b.WriteString(requestPayload)
		if err := adapter.ServeWebSocket(nil, b.Bytes()); err != nil {
			as.NoError(err)
			return
		}
		as.Equal(sum, 2)
	})
}