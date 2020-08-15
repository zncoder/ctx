package ctx

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/zncoder/assert"
	"golang.org/x/net/trace"
)

// Context adds tracing to a context.
// It encapsulates CancelFunc and provides Close to cancel.
type Context struct {
	context.Context
}

func New(parent context.Context) Context {
	if parent == nil {
		parent = context.Background()
	}
	cc, cancel := context.WithCancel(parent)
	cx := Context{cc}
	return cx.WithValue(closeKey{}, cancel)
}

// WithTrace associates trace with cx. Name is "family/title".
func (cx Context) WithTrace(name string) Context {
	family, title := splitTrace(name)
	assert.OK(family != "" && title != "")
	tr := trace.New(family, title)
	return cx.WithValue(traceKey{}, tr).WithValue(closeKey{}, tr.Finish)
}

func (cx Context) WithTimeout(timeout time.Duration) Context {
	cc, cancel := context.WithTimeout(cx.Context, timeout)
	cx = Context{cc}
	return cx.WithValue(closeKey{}, cancel)
}

// WithLog prints trace to log.
// To help trace, a random tag is chosen and printed in each log line.
func (cx Context) WithLog() Context {
	tag := strconv.FormatInt(int64(rand.Intn(60466176)), 36)
	return cx.WithValue(logKey{}, tag)
}

func (cx Context) WithValue(key, value interface{}) Context {
	return Context{context.WithValue(cx.Context, key, value)}
}

func (cx Context) Close() {
	if cancel, ok := cx.Value(closeKey{}).(context.CancelFunc); ok {
		cancel()
	}
}

func (cx Context) Printf(format string, args ...interface{}) {
	cx.PrintSkip(2, fmt.Sprintf(format, args...))
}

func (cx Context) Errorf(err error, format string, args ...interface{}) error {
	err = fmt.Errorf(format+" [err:%w]", append(args, err)...)
	cx.PrintSkip(2, err.Error())
	return err
}

func (cx Context) PrintSkip(skip int, s string) {
	if tag, ok := cx.Value(logKey{}).(string); ok {
		log.Output(skip+1, tag+": "+s)
	}

	if tr, ok := cx.Value(traceKey{}).(trace.Trace); ok {
		tr.LazyLog(stringer(s), false)
	}
}

func (cx Context) Sleep(min, max time.Duration) {
	assert.OK(min >= 0 && min <= max)

	delay := max - min
	if delay > 0 {
		delay = time.Duration(rand.Int63n(delay.Nanoseconds()))
	}
	delay += min

	tm := time.NewTimer(delay)
	select {
	case <-tm.C:
	case <-cx.Done():
	}
	tm.Stop()
}

type stringer string

func (s stringer) String() string { return string(s) }

func splitTrace(s string) (family, title string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ""
	}

	ss := strings.SplitN(s, "/", 2)
	if len(ss) == 1 {
		return "", ""
	}
	return ss[0], ss[1]
}

type (
	closeKey struct{}
	traceKey struct{}
	logKey   struct{}
)
