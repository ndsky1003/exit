package exit

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync/atomic"
	"syscall"
)

type iExitable interface {
	Exit() error
}

type call_back_func func() error

func (this call_back_func) Exit() error {
	return this()
}

type namefn struct {
	name     string
	fn       iExitable
	priority int
}

var arr []namefn

func Push(fn func() error, opts ...*options) {
	opt := NewOptions().Merge(opts...)
	arr = append(arr, namefn{name: opt.getName(), fn: call_back_func(fn), priority: opt.getPriority()})
	sort.Slice(arr, func(i, j int) bool { return arr[i].priority > arr[j].priority })
}

func PushExitable(fn iExitable, opts ...*options) {
	opt := NewOptions().Merge(opts...)
	arr = append(arr, namefn{name: opt.getName(), fn: fn, priority: opt.getPriority()})
	sort.Slice(arr, func(i, j int) bool { return arr[i].priority > arr[j].priority })
}

var exiting int32 = 0

func exit() {

	if b := atomic.CompareAndSwapInt32(&exiting, 0, 1); !b {
		return
	}

	fmt.Println("[EXIT] ################ 程序退出开始 ##############")

	for i := 0; i < len(arr); i++ {
		obj := arr[i]
		if err := obj.fn.Exit(); err != nil {
			fmt.Printf("[EXIT] ################ %v 执行失败:%v ##############\n", obj.name, err)
		}
	}

	fmt.Println("[EXIT] ################ 程序退出结束 ##############")

	os.Exit(0)
}

var ch chan os.Signal

func init() {

	ch = make(chan os.Signal)

	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ch
		exit()
	}()
}
