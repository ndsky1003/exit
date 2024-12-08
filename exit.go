package exit

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type IExitable interface {
	Exit() error
}

type call_back_func func() error

func (this call_back_func) Exit() error {
	return this()
}

type nameFn struct {
	fn  IExitable
	opt *Option
}

type ExitManager struct {
	tasks   []nameFn
	exiting int32
	ch      chan os.Signal
	l       sync.Mutex
	opt     *OptionMgr
}

func NewExitManager(opts ...*OptionMgr) *ExitManager {
	opt := OptionsMgr().
		SetLockFile(fmt.Sprintf("/tmp/%s.lock", filepath.Base(os.Args[0]))).
		Merge(opts...)

	lockFile := *opt.lockedfile
	if checkLock(lockFile) {
		fmt.Println("Another instance is already running. Exiting...")
		os.Exit(1)
	}

	if err := createLock(lockFile); err != nil {
		fmt.Printf("Failed to create lock file: %v\n", err)
		os.Exit(1)
	}
	em := &ExitManager{
		ch:  make(chan os.Signal),
		opt: opt,
	}
	signal.Notify(em.ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go em.waitForSignal()
	return em
}

func (this *ExitManager) Push(fn func() error, opts ...*Option) {
	opt := Options().
		SetPriority(0).
		SetTimeout(10 * time.Second).
		SetName("").
		Merge(opts...)
	this.l.Lock()
	defer this.l.Unlock()
	this.tasks = append(this.tasks, nameFn{fn: call_back_func(fn), opt: opt})
	sort.Slice(this.tasks, func(i, j int) bool { return *this.tasks[i].opt.priority > *this.tasks[j].opt.priority })
}

func (this *ExitManager) PushExitable(fn IExitable, opts ...*Option) {
	opt := Options().Merge(opts...)
	this.l.Lock()
	defer this.l.Unlock()
	this.tasks = append(this.tasks, nameFn{opt: opt, fn: fn})
	sort.Slice(this.tasks, func(i, j int) bool { return *this.tasks[i].opt.priority > *this.tasks[j].opt.priority })
}

func (this *ExitManager) exit() {
	if b := atomic.CompareAndSwapInt32(&this.exiting, 0, 1); !b {
		return
	}

	fmt.Println("[EXIT] ################ 程序退出开始 ##############")

	this.l.Lock()
	defer this.l.Unlock()
	lockfile := *this.opt.lockedfile

	for _, obj := range this.tasks {
		done := make(chan error)
		go func(f IExitable) {
			done <- f.Exit()
		}(obj.fn)
		timeout := *obj.opt.timeout
		name := *obj.opt.name

		select {
		case err := <-done:
			if err != nil {
				fmt.Printf("[EXIT] ################ %v 执行失败:%v ##############\n", name, err)
			} else {
				fmt.Printf("[EXIT] ################ %v 执行成功 ##############\n", name)
			}
		case <-time.After(timeout):
			fmt.Printf("[EXIT] ################ %v 执行超时 ##############\n", name)
		}
	}

	fmt.Println("[EXIT] ################ 程序退出结束 ##############")
	os.Remove(lockfile) // 删除锁文件
	os.Exit(0)
}

func (em *ExitManager) waitForSignal() {
	<-em.ch
	em.exit()
}

func checkLock(lockFile string) bool {
	if _, err := os.Stat(lockFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func createLock(lockFile string) error {
	file, err := os.Create(lockFile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
