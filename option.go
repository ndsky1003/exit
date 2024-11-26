package exit

import "time"

type Option struct {
	name     *string
	priority *int
	timeout  *time.Duration
}

func Options() *Option {
	return &Option{}
}

func (this *Option) SetPriority(p int) *Option {
	this.priority = &p
	return this
}
func (this *Option) getPriority() int {
	if this == nil || this.priority == nil {
		return 0
	}
	return *this.priority
}

func (this *Option) SetName(n string) *Option {
	this.name = &n
	return this
}

func (this *Option) getName() string {
	if this == nil || this.name == nil {
		return ""
	}
	return *this.name
}

func (this *Option) SetTimeout(timeout time.Duration) *Option {
	this.timeout = &timeout
	return this
}

func (this *Option) merge(opt *Option) {
	if opt == nil {
		return
	}
	if opt.priority != nil {
		this.priority = opt.priority
	}
	if opt.name != nil {
		this.name = opt.name
	}
	if opt.timeout != nil {
		this.timeout = opt.timeout
	}
}

func (this *Option) Merge(opts ...*Option) *Option {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}

type OptionMgr struct {
	lockedfile *string
}

func OptionsMgr() *OptionMgr {
	return &OptionMgr{}
}

func (this *OptionMgr) SetLockFile(file string) *OptionMgr {
	this.lockedfile = &file
	return this
}

func (this *OptionMgr) Merge(opts ...*OptionMgr) *OptionMgr {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}

func (this *OptionMgr) merge(opt *OptionMgr) {
	if opt == nil {
		return
	}
	if opt.lockedfile != nil {
		this.lockedfile = opt.lockedfile
	}
}
