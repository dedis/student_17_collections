package collection

import "testing"

type testctx struct {
    file string
    test *testing.T
}

func (this testctx) should_panic(prefix string, function func()) {
    defer func() {
        if recover() == nil {
            this.test.Error(this.file, prefix, "Function provided did not panic.")
        }
    }()

    function()
}
