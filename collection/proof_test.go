package collection

import "testing"

func TestProofGetters(test *testing.T) {
    key := []byte("mykey")
    value := []byte("myvalue")

    left := dump{}
    right := dump{}

    left.key = []byte("wrongkey")
    right.key = key
    right.values = [][]byte{value}

    steps := []step{{}, {}, {}, {}, {left, right}}

    proof := Proof{}
    proof.key = key
    proof.steps = steps

    if string(proof.Key()) != string(key) {
        test.Error("[key]", "Key returns wrong key.")
    }

    if len(proof.Values()) != 1 {
        test.Error("[values]", "Values returns wrong number of values.")
    }

    if string(proof.Values()[0]) != string(value) {
        test.Error("[values]", "Values returns wrong values.")
    }

    if !(proof.Match()) {
        test.Error("[match]", "Match is false on matching proof.")
    }


    nomatch := dump{}
    nomatch.key = []byte("nomatch")
    nomatch.values = [][]byte{[]byte("these"), []byte("values"), []byte("don't"), []byte("concern"), []byte("you")}

    steps = []step{{}, {}, {}, {}, {left, nomatch}}
    proof.steps = steps

    if string(proof.Key()) != string(key) {
        test.Error("[key]", "Key returns wrong key.")
    }

    if len(proof.Values()) != 0 {
        test.Error("[values]", "Values returns values on non-matching proof.")
    }

    if proof.Match() {
        test.Error("[match]", "Match is true on non-matching proof.")
    }

    proof.steps = []step{}

    if len(proof.Values()) != 0 {
        test.Error("[values]", "Values returns values on empty proof.")
    }

    if proof.Match() {
        test.Error("[match]", "Match is true on empty proof.")
    }
}
