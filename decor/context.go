package decor

type TKind int

const (
	KFunc TKind = iota
)

type Context struct {
	Kind TKind
	TargetIn,
	TargetOut []any
	Func func()

	// inner
	doRef int64
}

func (d *Context) TargetDo() {
	d.doRef++
	d.Func()
}
