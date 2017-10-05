package collection

import "testing"
import "math/rand"
import "encoding/binary"

// Helper

type collectiontesthelper struct {
}

func (this collectiontesthelper) validatenode(collection * collection, node * node) bool {
    if !(node.known) {
        return true
    }

    if node.leaf() {
        if (node.children.left != nil) || (node.children.right != nil) {
            return false
        }

        if node.label != sha256(true, node.key, node.values) {
            return false
        }
    } else {
        if (node.children.left == nil) || (node.children.right == nil) {
            return false
        }

        if (node.children.left.parent != node) || (node.children.right.parent != node) {
            return false
        }

        if node.label != sha256(false, node.values, node.children.left.label[:], node.children.right.label[:]) {
            return false
        }

        if node.children.left.known && node.children.right.known {
            for index := 0; index < len(collection.values); index++ {
                value := collection.values[index].Parent(node.children.left.values[index], node.children.right.values[index])
                if (len(value) != len(node.values[index])) || !match(value, node.values[index], 8 * len(value)) {
                    return false
                }
            }
        }
    }

    return true
}

func (this collectiontesthelper) validatetree(collection * collection, node * node, path []bool) bool {
    if !(this.validatenode(collection, node)) {
        return false
    }

    if node.leaf() {
        if !(node.placeholder()) {
            for index := 0; index < len(path); index++ {
                keyhash := sha256(node.key)
                if path[index] != bit(keyhash[:], index) {
                    return false
                }
            }
        }

        return true
    } else {
        leftpath := make([]bool, len(path))
        rightpath := make([]bool, len(path))

        copy(leftpath, path)
        copy(rightpath, path)

        leftpath = append(leftpath, false)
        rightpath = append(rightpath, true)

        return (this.validatetree(collection, node.children.left, leftpath) && this.validatetree(collection, node.children.right, rightpath))
    }
}

// Tests

func TestEmptyCollection(test *testing.T) {
    helper := collectiontesthelper{}

    basecollection := EmptyCollection()

    if !(basecollection.root.known) || !(basecollection.root.children.left.known) || !(basecollection.root.children.right.known) {
        test.Error("[known]", "New collection has unknown nodes.")
    }

    if !(basecollection.root.root()) {
        test.Error("[root]", "Collection root is not a root.")
    }

    if basecollection.root.leaf() {
        test.Error("[root]", "Collection root doesn't have children.")
    }

    if !(basecollection.root.children.left.placeholder()) || !(basecollection.root.children.right.placeholder()) {
        test.Error("[leaves]", "Collection leaves are not placeholder leaves.")
    }

    if len(basecollection.root.values) != 0 || len(basecollection.root.children.left.values) != 0 || len(basecollection.root.children.right.values) != 0 {
        test.Error("[values]", "Nodes of a collection without values have values.")
    }

    if !(helper.validatetree(&basecollection, basecollection.root, []bool{})) {
        test.Error("[tree]", "Invalid collection tree.")
    }

    stake64 := Stake64{}
    stakecollection := EmptyCollection(stake64)

    if len(stakecollection.root.values) != 1 || len(stakecollection.root.children.left.values) != 1 || len(stakecollection.root.children.right.values) != 1 {
        test.Error("[values]", "Nodes of a stake collection don't have exactly one value.")
    }

    if stake64.Decode(stakecollection.root.values[0]) != 0 || stake64.Decode(stakecollection.root.children.left.values[0]) != 0 || stake64.Decode(stakecollection.root.children.right.values[0]) != 0 {
        test.Error("[stake]", "Nodes of an empty stake collection don't have zero stake.")
    }

    if !(helper.validatetree(&stakecollection, stakecollection.root, []bool{})) {
        test.Error("[tree]", "Invalid stake collection tree.")
    }

    data := Data{}
    stakedatacollection := EmptyCollection(stake64, data)

    if len(stakedatacollection.root.values) != 2 || len(stakedatacollection.root.children.left.values) != 2 || len(stakedatacollection.root.children.right.values) != 2 {
        test.Error("[values]", "Nodes of a data and stake collection don't have exactly one value.")
    }

    if len(stakedatacollection.root.values[1]) != 0 || len(stakedatacollection.root.children.left.values[1]) != 0 || len(stakedatacollection.root.children.right.values[1]) != 0 {
        test.Error("[values]", "Nodes of a data and stake collection don't have empty data value.")
    }

    if !(helper.validatetree(&stakedatacollection, stakedatacollection.root, []bool{})) {
        test.Error("[tree]", "Invalid data and stake collection tree.")
    }
}

func TestBegin(test *testing.T) {
    collection := EmptyCollection()
    collection.Begin()

    if !(collection.transaction) {
        test.Error("[transaction]", "Begin does not set transaction flag to true.")
    }
}

func TestEnd(test *testing.T) {
    helper := collectiontesthelper{}
    collection := EmptyCollection()

    collection.Begin()

    for index := uint64(0); index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, index)

        collection.Apply(AddUpdate(key, [][]byte{}))
    }

    collection.End()

    if collection.transaction {
        test.Error("[transaction]", "End does not set transaction flag to false.")
    }

    if !(helper.validatetree(&collection, collection.root, []bool{})) {
        test.Error("[tree]", "Ending a transaction after multiple adds produces an invalid tree.")
        return
    }

    noscope := EmptyCollection()
    noscope.Scope.None()

    noscope.Begin()

    for index := uint64(0); index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, index)

        noscope.Apply(AddUpdate(key, [][]byte{}))
    }

    noscope.End()

    if noscope.root.label != collection.root.label {
        test.Error("[noscope]", "Ending the same transaction on a collection with empty scope produces a different root label.")
    }

    if noscope.root.children.left.children.left.known || noscope.root.children.left.children.right.known || noscope.root.children.right.children.left.known || noscope.root.children.right.children.right.known {
        test.Error("[pruning]", "A collection with empty scope should not have known root grandchildren.")
    }
}

func TestPlaceholderValues(test *testing.T) {
    data := Data{}
    stake64 := Stake64{}

    basecollection := EmptyCollection()
    values := basecollection.placeholdervalues()

    if len(values) != 0 {
        test.Error("[base]", "Base collection yields non-empty list of placeholder values.")
    }

    datacollection := EmptyCollection(data)
    values = datacollection.placeholdervalues()

    if len(values) != 1 {
        test.Error("[data]", "Data collection does not yield list of placeholder values with exactly one element.")
    }

    if len(values[0]) != 0 {
        test.Error("[data]", "Data collection yields non-empty placeholder value.")
    }

    stakecollection := EmptyCollection(stake64)
    values = stakecollection.placeholdervalues()

    if len(values) != 1 {
        test.Error("[stake]", "Stake collection does not yield list of placeholder values with exactly one element.")
    }

    if stake64.Decode(values[0]) != 0 {
        test.Error("[stake]", "Stake collection yields non-zero placeholder value.")
    }

    multicollection := EmptyCollection(Data{}, Stake64{}, Stake64{}, Data{}, Data{}, Stake64{}, Data{})
    values = multicollection.placeholdervalues()

    if len(values) != 7 {
        test.Error("[multi]", "Multi collection yields wrong number of placeholder values.")
    }

    if (len(values[0]) != 0) || (stake64.Decode(values[1]) != 0) || (stake64.Decode(values[2]) != 0) || (len(values[3]) != 0) || (len(values[4]) != 0) || (stake64.Decode(values[5]) != 0) || (len(values[6]) != 0) {
        test.Error("[multi]", "Multi collection yields wrong placeholder values.")
    }
}

func TestUpdate(test *testing.T) {
    helper := collectiontesthelper{}

    basecollection := EmptyCollection()

    error := basecollection.update(&node{})
    if error == nil {
        test.Error("[known]", "Update doesn't yield error on unknown node.")
    }

    basecollection.root.children.left.known = false
    error = basecollection.update(basecollection.root)

    if error == nil {
        test.Error("[known]", "Update doesn't yield error on node with unknown children.")
    }

    stake64 := Stake64{}
    stakecollection := EmptyCollection(stake64)

    stakecollection.root.children.left.values[0] = stake64.Encode(66)

    if stakecollection.update(stakecollection.root.children.left) != nil {
        test.Error("[stake]", "Update fails on stake leaf.")
    }

    if stakecollection.update(stakecollection.root) != nil {
        test.Error("[stake]", "Update fails on stake root.")
    }

    if stake64.Decode(stakecollection.root.values[0]) != 66 {
        test.Error("[stake]", "Wrong value on stake root.")
    }

    stakecollection.root.children.right.values[0] = stake64.Encode(33)

    if stakecollection.update(stakecollection.root.children.right) != nil {
        test.Error("[stake]", "Update fails on stake leaf.")
    }

    if stakecollection.update(stakecollection.root) != nil {
        test.Error("[stake]", "Update fails on stake root.")
    }

    if stake64.Decode(stakecollection.root.values[0]) != 99 {
        test.Error("[stake]", "Wrong value on stake root.")
    }

    if !(helper.validatetree(&stakecollection, stakecollection.root, []bool{})) {
        test.Error("[tree]", "Invalid stake collection tree.")
    }
}

func TestFix(test *testing.T) {
    helper := collectiontesthelper{}

    collection := EmptyCollection()

    collection.root.children.left.key = []byte("leftkey")
    collection.root.children.left.inconsistent = true
    collection.root.inconsistent = true

    if helper.validatetree(&collection, collection.root, []bool{}) {
        test.Error("[inconsistent]", "Inconsistent tree should not be valid.")
    }

    collection.fix(collection.root)

    if !(helper.validatetree(&collection, collection.root, []bool{})) {
        test.Error("[fix]", "Tree should be valid after fix.")
    }

    oldrootlabel := collection.root.label

    collection.root.children.right.key = []byte("rightkey")
    collection.root.children.right.inconsistent = true

    collection.fix(collection.root)

    if collection.root.label != oldrootlabel {
        test.Error("[fix]", "Fix should not visit nodes that are not marked as inconsistent.")
    }

    collection.root.inconsistent = true
    collection.fix(collection.root)

    if collection.root.label == oldrootlabel {
        test.Error("[fix]", "Fix should alter the label of the root of a collection tree.")
    }

    if !(helper.validatetree(&collection, collection.root, []bool{})) {
        test.Error("[fix]", "Tree should be valid after fix.")
    }
}

func TestCollect(test *testing.T) {
    collection := EmptyCollection()

    collection.root.children.left.children.left = new(node)
    collection.root.children.left.children.right = new(node)

    collection.root.children.right.children.left = new(node)
    collection.root.children.right.children.left.children.left = new(node)
    collection.root.children.right.children.left.children.right = new(node)
    collection.root.children.right.children.right = new(node)
    collection.root.children.right.children.right.children.left = new(node)
    collection.root.children.right.children.right.children.right = new(node)

    collection.temporary = append(collection.temporary, collection.root.children.left, collection.root.children.right)
    collection.collect()

    if collection.root.children.left.known || collection.root.children.right.known {
        test.Error("[known]", "Collect doesn't make temporary nodes unknown.")
    }

    if (collection.root.children.left.children.left != nil) || (collection.root.children.left.children.right != nil) || (collection.root.children.right.children.left != nil) || (collection.root.children.right.children.right != nil) {
        test.Error("[children]", "Collect does not prune children of temporary nodes.")
    }

    if len(collection.temporary) != 0 {
        test.Error("[temporary]", "Collect does not empty temporary nodes list.")
    }
}

func TestApplyAdd(test *testing.T) {
    helper := collectiontesthelper{}

    stake64 := Stake64{}
    collection := EmptyCollection(stake64)

    for index := uint64(0); index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, index)

        collection.Apply(AddUpdate(key, [][]byte{stake64.Encode(uint64(rand.Uint32()))}))

        if !(helper.validatetree(&collection, collection.root, []bool{})) {
            test.Error("[tree]", "Add produces an invalid tree.")
            return
        }
    }

    unknownroot := EmptyCollection()
    unknownroot.root.known = false

    error := unknownroot.Apply(AddUpdate([]byte("key"), [][]byte{}))
    if error == nil {
        test.Error("[unknownroot]", "Add should yield an error on a collection with unknown root.")
    }

    unknownrootchildren := EmptyCollection()
    unknownrootchildren.root.children.left.known = false
    unknownrootchildren.root.children.right.known = false

    error = unknownrootchildren.Apply(AddUpdate([]byte("key"), [][]byte{}))
    if error == nil {
        test.Error("[unknownrootchildren]", "Add should yield an error on a collection with unknown root children.")
    }

    keycollision := EmptyCollection()
    keycollision.Apply(AddUpdate([]byte("key"), [][]byte{}))

    error = keycollision.Apply(AddUpdate([]byte("key"), [][]byte{}))
    if error == nil {
        test.Error("[keycollision]", "Add should yield an error on key collision.")
    }
}
