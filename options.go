package exit

type options struct {
	name     *string
	priority *int
}

func NewOptions() *options {
	return &options{}
}

func (this *options) SetPriority(p int) *options {
	this.priority = &p
	return this
}
func (this *options) getPriority() int {
	if this == nil || this.priority == nil {
		return 0
	}
	return *this.priority
}

func (this *options) SetName(n string) *options {
	this.name = &n
	return this
}

func (this *options) getName() string {
	if this == nil || this.name == nil {
		return ""
	}
	return *this.name
}

func (this *options) merge(opt *options) {
	if opt.priority != nil {
		this.priority = opt.priority
	}
	if opt.name != nil {
		this.name = opt.name
	}
}

func (this *options) Merge(opts ...*options) *options {
	for _, opt := range opts {
		this.merge(opt)
	}
	return this
}
