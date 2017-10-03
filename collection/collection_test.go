package collection

import "testing"

func TestEmptyCollection(test *testing.T) {
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

    stake64 := Stake64{}
    stakecollection := EmptyCollection(stake64)

    if len(stakecollection.root.values) != 1 || len(stakecollection.root.children.left.values) != 1 || len(stakecollection.root.children.right.values) != 1 {
        test.Error("[values]", "Nodes of a stake collection don't have exactly one value.")
    }

    if stake64.Decode(stakecollection.root.values[0]) != 0 || stake64.Decode(stakecollection.root.children.left.values[0]) != 0 || stake64.Decode(stakecollection.root.children.right.values[0]) != 0 {
        test.Error("[stake]", "Nodes of an empty stake collection don't have zero stake.")
    }

    data := Data{}
    stakedatacollection := EmptyCollection(stake64, data)

    if len(stakedatacollection.root.values) != 2 || len(stakedatacollection.root.children.left.values) != 2 || len(stakedatacollection.root.children.right.values) != 2 {
        test.Error("[values]", "Nodes of a data and stake collection don't have exactly one value.")
    }

    if len(stakedatacollection.root.values[1]) != 0 || len(stakedatacollection.root.children.left.values[1]) != 0 || len(stakedatacollection.root.children.right.values[1]) != 0 {
        test.Error("[values]", "Nodes of a data and stake collection don't have empty data value.")
    }
}

func TestUpdate(test *testing.T) {
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

    // TODO: Check for SHA256 consistency with verification procedures (yet to be implemented).
}
