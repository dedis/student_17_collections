package collection

import "testing"

func TestRecord(test *testing.T) {
    stake64 := Stake64{}
    data := Data{}

    collection := EmptyCollection(stake64, data)
    collection.Add([]byte("mykey"), uint64(66), []byte("mydata"))

    var leaf *node

    if collection.root.children.left.placeholder() {
        leaf = collection.root.children.right
    } else {
        leaf = collection.root.children.left
    }

    match := recordmatch(&collection, leaf)
    mismatch := recordmismatch(&collection, []byte("wrongkey"))

    if (match.collection != &collection) || (mismatch.collection != &collection) {
        test.Error("[record.go]", "[constructors]", "Constructors don't set collection appropriately.")
    }

    if !(match.match) || mismatch.match {
        test.Error("[record.go]", "[constructors]", "Constructors don't set match appropriately.")
    }

    if !equal(match.key, []byte("mykey")) || !equal(mismatch.key, []byte("wrongkey")) {
        test.Error("[record.go]", "[constructors]", "Constructors don't set key appropriately")
    }

    if len(match.values) != 2 || len(mismatch.values) != 0 {
        test.Error("[record.go]", "[constructors]", "Constructors don't set the appropriate number of values.")
    }

    if !equal(match.values[0], leaf.values[0]) || !equal(match.values[1], leaf.values[1]) {
        test.Error("[record.go]", "[constructors]", "Constructors set the wrong values.")
    }

    if !(match.Match()) || mismatch.Match() {
        test.Error("[record.go]", "[match]", "Match() returns the wrong value.")
    }

    if !equal(match.Key(), []byte("mykey")) || !equal(mismatch.Key(), []byte("wrongkey")) {
        test.Error("[record.go]", "[key]", "Key() returns the wrong value.")
    }

    matchvalues, matcherror := match.Values()

    if matcherror != nil {
        test.Error("[record.go]", "[values]", "Values() yields an error on matching record.")
    }

    if len(matchvalues) != 2 {
        test.Error("[record.go]", "[values]", "Values() returns the wrong number of values")
    }

    if (matchvalues[0].(uint64) != 66) || !equal(matchvalues[1].([]byte), leaf.values[1]) {
        test.Error("[record.go]", "[values]", "Values() returns the wrong values.")
    }

    _, mismatcherror := mismatch.Values()

    if mismatcherror == nil {
        test.Error("[record.go]", "[values]", "Values() does not yield an error on mismatching record.")
    }

    match.values[0] = match.values[0][:6]

    _, illformederror := match.Values()

    if illformederror == nil {
        test.Error("[record.go]", "[values]", "Values() does not yield an error on record with ill-formed values.")
    }

    match.values = match.values[:1]

    _, fewerror := match.Values()

    if fewerror == nil {
        test.Error("[record.go]", "[values]", "Values() does not yield an error on record with wrong number of values.")
    }
}
