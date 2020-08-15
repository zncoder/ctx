package ctx_test

import (
	"io"
	"log"
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
	_ = cx.Errorf(io.EOF, "EOF error")
	cx.Close()

	type key1 struct{}
	type key2 struct{}
	cx1 := ctx.New(nil).WithValue(key1{}, "key1")
	cx2 := cx1.WithValue(key2{}, "key2")
	if v1, ok := cx1.Value(key1{}).(string); ok {
		if v1 != "key1" {
			log.Fatal(v1)
		}
	} else {
		log.Fatal("key1 not found")
	}
	if _, ok := cx1.Value(key2{}).(string); ok {
		log.Fatal("key2 shouldn't be found")
	}
	if v1, ok := cx2.Value(key1{}).(string); ok {
		if v1 != "key1" {
			log.Fatal(v1)
		}
	} else {
		log.Fatal("key1 not found")
	}
	if v2, ok := cx2.Value(key2{}).(string); ok {
		if v2 != "key2" {
			log.Fatal(v2)
		}
	} else {
		log.Fatal("key2 not found")
	}

	// Output:
}
