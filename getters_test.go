package collection

import "testing"
import "encoding/binary"

func TestGettersConstructors(test *testing.T) {
    collection := EmptyCollection()
    getter := collection.Get([]byte("mykey"))

    if getter.collection != &collection {
        test.Error("[getters.go]", "[constructors]", "Getter constructor sets wrong collection pointer.")
    }

    if !equal(getter.key, []byte("mykey")) {
        test.Error("[getters.go]", "[constructors]", "Getter constructor sets wrong key.")
    }
}

func TestGettersRecord(test *testing.T) {
    collection := EmptyCollection()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
    }

    for index := 0; index < 1024; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        record, err := collection.Get(key).Record()

        if err != nil {
            test.Error("[getters.go]", "[record]", "Record() yields an error on valid key query.")
        }

        if !equal(record.Key(), key) {
            test.Error("[getters.go]", "[record]", "Record() returns a record with wrong key.")
        }

        if (index < 512) && !(record.Match()) {
            test.Error("[getters.go]", "[record]", "Record() yields a non-matching record on existing key.")
        }

        if (index >= 512) && record.Match() {
            test.Error("[getters.go]", "[record]", "Record() yields a matching record on non-existing key.")
        }
    }

    collection.Scope.None()
    collection.collect()

    _, err := collection.Get([]byte("mykey")).Record()

    if err == nil {
        test.Error("[getters.go]", "[record]", "Record() does not yield an error on unknown collection.")
    }
}

func TestGettersProof(test *testing.T) {
    collection := EmptyCollection()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key)
    }

    for index := 0; index < 1024; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        proof, err := collection.Get(key).Proof()

        if err != nil {
            test.Error("[getters.go]", "[proof]", "Proof() yields an error on valid key query.")
        }

        if !equal(proof.Key(), key) {
            test.Error("[getters.go]", "[proof]", "Proof() returns a record with wrong key.")
        }

        if (index < 512) && !(proof.Match()) {
            test.Error("[getters.go]", "[proof]", "Proof() yields a non-matching record on existing key.")
        }

        if (index >= 512) && proof.Match() {
            test.Error("[getters.go]", "[proof]", "Proof() yields a matching record on non-existing key.")
        }

        if proof.collection != &collection {
            test.Error("[getters.go]", "[proof]", "Proof() returns proof with wrong collection pointer.")
        }

        if proof.root.label != collection.root.label {
            test.Error("[getters.go]", "[proof]", "Proof() returns a proof with wrong root.")
        }

        if !(proof.root.consistent()) {
            test.Error("[getters.go]", "[proof]", "Proof() returns a proof with inconsistent root.")
        }

        if len(proof.steps) == 0 {
            test.Error("[getters.go]", "[proof]", "Proof() returns a proof with no steps.")
        }

        if (proof.steps[0].left.label != proof.root.children.left) || (proof.steps[0].right.label != proof.root.children.right) {
            test.Error("[getters.go]", "[proof]", "Label mismatch between root children and first step.")
        }

        path := sha256(key)

        for depth := 0; depth < len(proof.steps) - 1; depth++ {
            if !(proof.steps[depth].left.consistent()) || !(proof.steps[depth].right.consistent()) {
                test.Error("[getters.go]", "[proof]", "Inconsistent step.")
            }

            if bit(path[:], depth) {
                if (proof.steps[depth].right.children.left != proof.steps[depth + 1].left.label) || (proof.steps[depth].right.children.right != proof.steps[depth + 1].right.label) {
                    test.Error("[getters.go]", "[proof]", "Step label mismatch given path.")
                }
            } else {
                if (proof.steps[depth].left.children.left != proof.steps[depth + 1].left.label) || (proof.steps[depth].left.children.right != proof.steps[depth + 1].right.label) {
                    test.Error("[getters.go]", "[proof]", "Step label mismatch given path.")
                }
            }
        }

        if !(proof.steps[len(proof.steps) - 1].left.consistent()) || !(proof.steps[len(proof.steps) - 1].right.consistent()) {
            test.Error("[getters.go]", "[proof]", "Last inconsistent step.")
        }
    }

    collection.Scope.Add([]byte{0xff}, 1)
    collection.collect()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        path := sha256(key)

        if bit(path[:], 0) {
            continue
        }

        _, err := collection.Get(key).Proof()
        if err == nil {
            test.Error("[getters.go]", "[proof]", "Proof() does not yield an error when querying an unknown subtree.")
        }
    }

    collection.Scope.None()
    collection.collect()

    _, err := collection.Get([]byte("mykey")).Proof()
    if err == nil {
        test.Error("[getters.go]", "[proof]", "Proof() does not yield an error when querying a tree with unknown root.")
    }
}
