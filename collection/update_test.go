package collection

import "testing"

func TestUpdateConstructors(test *testing.T) {
    addkey := []byte("keytoadd")
    addvalues := [][]byte{[]byte("One pretty value")}
    addreq := AddUpdate(addkey, addvalues)

    if addreq.kind != add {
        test.Error("[add]", "Wrong update kind.")
    }

    if (len(addreq.key) != len(addkey)) || (!match(addreq.key, addkey, 8 * len(addkey))) {
        test.Error("[add]", "Wrong key.")
    }

    if (len(addreq.values) != 1) || (len(addreq.values[0]) != len(addvalues[0])) || (!match(addreq.values[0], addvalues[0], 8 * len(addvalues[0]))) {
        test.Error("[add]", "Wrong values.")
    }

    removekey := []byte("keytoremove")
    removereq := RemoveUpdate(removekey)

    if removereq.kind != remove {
        test.Error("[remove]", "Wrong update kind.")
    }

    if (len(removereq.key) != len(removekey)) || (!match(removereq.key, removekey, 8 * len(removekey))) {
        test.Error("[remove]", "Wrong key.")
    }

    if len(removereq.values) != 0 {
        test.Error("[remove]", "Values in remove update.")
    }

    setkey := []byte("keytoset")
    setvalues := [][]byte{[]byte("Another pretty value")}
    setreq := SetUpdate(setkey, setvalues)

    if setreq.kind != set {
        test.Error("[set]", "Wrong update kind.")
    }

    if (len(setreq.key) != len(setkey)) || (!match(setreq.key, setkey, 8 * len(setkey))) {
        test.Error("[set]", "Wrong key.")
    }

    if (len(setreq.values) != 1) || (len(setreq.values[0]) != len(setvalues[0])) || (!match(setreq.values[0], setvalues[0], 8 * len(setvalues[0]))) {
        test.Error("[set]", "Wrong values.")
    }
}
