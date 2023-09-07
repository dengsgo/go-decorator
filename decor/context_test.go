package decor

import "testing"

func TestContext_DoRef(t *testing.T) {
	ctx := &Context{
		Func: func() {
			// whatever you do
		},
	}
	var i int64 = 1
	for ; i < 100; i++ {
		ctx.TargetDo()
		if ctx.DoRef() != i {
			t.Fatal("ctx.DoRef() != i, i=", i)
		}
	}
}

func TestContext_TargetDo(t *testing.T) {
	i := 100
	s := ""
	ctx := &Context{
		Func: func() {
			func() {
				i = 100 * 10
				s = "TargetDo()"
			}() // whatever you do
		},
	}
	ctx.TargetDo()
	if i != 1000 {
		t.Fatal("i want 1000, but get", i)
	}
	if s != "TargetDo()" {
		t.Fatal("s want `TargetDo()`, but get `", i, "`")
	}
}
