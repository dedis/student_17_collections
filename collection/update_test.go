package collection

import "testing"

func TestAddUpdate(test *testing.T) {
    collection := EmptyCollection()

    proof, _ := collection.Get([]byte("mykey")).Proof()
    update := AddUpdate([]byte("mykey"), [][]byte{}, proof)

    if !(collection.Verify(update)) {
        test.Error("[verify]", "Collection fails to verify valid update.")
    }

    if !(collection.Apply(update)) {
        test.Error("[apply]", "Apply returns false on applicable update.")
    }

    record, _ := collection.Get([]byte("mykey")).Record()

    if !(record.Exists()) {
        test.Error("[apply]", "Record is not successfully added when calling collection.Apply.")
    }

    collection = EmptyCollection()

    proof, _ = collection.Get([]byte("mykey")).Proof()
    update = AddUpdate([]byte("mykey"), [][]byte{}, proof)

    collection.Begin()

    if !(collection.Verify(update)) {
        test.Error("[verify]", "Collection fails to verify valid update.")
    }

    if !(collection.Verify(update)) {
        test.Error("[verify]", "Collection fails to verify valid update.")
    }

    if !(collection.Apply(update)) {
        test.Error("[apply]", "Apply returns false on applicable update.")
    }

    if collection.Apply(update) {
        test.Error("[apply]", "Apply returns true on non-applicable update (collision not detected).")
    }

    collection.End()

    record, _ = collection.Get([]byte("mykey")).Record()

    if !(record.Exists()) {
        test.Error("[apply]", "Record is not successfully added when calling collection.Apply.")
    }
}
