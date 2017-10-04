package collection

import "testing"

func TestRequestConstructors(test *testing.T) {
    addkey := []byte("keytoadd")
    addvalues := [][]byte{[]byte("One pretty value")}
    addreq := AddRequest(addkey, addvalues)

    if addreq.kind != add {
        test.Error("[add]", "Wrong request kind.")
    }

    if (len(addreq.key) != len(addkey)) || (!match(addreq.key, addkey, 8 * len(addkey))) {
        test.Error("[add]", "Wrong key.")
    }

    if (len(addreq.values) != 1) || (len(addreq.values[0]) != len(addvalues[0])) || (!match(addreq.values[0], addvalues[0], 8 * len(addvalues[0]))) {
        test.Error("[add]", "Wrong values.")
    }

    removekey := []byte("keytoremove")
    removereq := RemoveRequest(removekey)

    if removereq.kind != remove {
        test.Error("[remove]", "Wrong request kind.")
    }

    if (len(removereq.key) != len(removekey)) || (!match(removereq.key, removekey, 8 * len(removekey))) {
        test.Error("[remove]", "Wrong key.")
    }

    if len(removereq.values) != 0 {
        test.Error("[remove]", "Values in remove request.")
    }

    updatekey := []byte("keytoupdate")
    updatevalues := [][]byte{[]byte("Another pretty value")}
    updatereq := UpdateRequest(updatekey, updatevalues)

    if updatereq.kind != update {
        test.Error("[update]", "Wrong request kind.")
    }

    if (len(updatereq.key) != len(updatekey)) || (!match(updatereq.key, updatekey, 8 * len(updatekey))) {
        test.Error("[update]", "Wrong key.")
    }

    if (len(updatereq.values) != 1) || (len(updatereq.values[0]) != len(updatevalues[0])) || (!match(updatereq.values[0], updatevalues[0], 8 * len(updatevalues[0]))) {
        test.Error("[update]", "Wrong values.")
    }
}
