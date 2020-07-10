package ctx_test

import (
	"io"
	"time"

	"github.com/zncoder/ctx"
)

func Example() {
	cx := ctx.New(nil)
	cx.Close()

	cx = ctx.New(nil).WithTrace("foo/bar")
	cx.Printf("trace an operation")
	cx.Close()

	cx = ctx.New(nil).WithTimeout(time.Millisecond)
	cx.Sleep(time.Hour, time.Hour)
	cx.Close()

	cx = ctx.New(nil).WithLog()
	cx.Printf("log an operation")
	cx.Close()

	cx = ctx.New(nil).WithLog()
	_ = cx.Errorf("err is %w", io.EOF)
	cx.Close()

	// Output:
}
