package collection

import "testing"

func TestRecordGetters(test *testing.T) {
    empty := record{}

    if empty.Exists() {
        test.Error("[exists]", "Exists returns true on empty record.")
    }

    if len(empty.Key()) != 0 {
        test.Error("[key]", "Key returns non-null key on empty record.")
    }

    if len(empty.Values()) != 0 {
        test.Error("[values]", "Values returns non-null values on empty record.")
    }

    something := record{nil, []byte("Hello World"), [][]byte{[]byte("Hello"), []byte("World")}}

    if !(something.Exists()) {
        test.Error("[exists]", "Exists returns false on non-empty record.")
    }

    if string(something.Key()) != "Hello World" {
        test.Error("[key]", "Key returns wrong key on non-empty record.")
    }

    if (len(something.Values()) != 2) || (string(something.Values()[0]) != "Hello") || (string(something.Values()[1]) != "World") {
        test.Error("[values]", "Values returns wrong values on non-empty record.")
    }
}
