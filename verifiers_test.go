package collection

import "testing"
import "encoding/binary"

func TestVerifiersVerify(test *testing.T) {
    ctx := testctx("[verifiers.go]", test)

    stake64 := Stake64{}
    data := Data{}

    collection := EmptyCollection(stake64, data)
    unknown := EmptyCollection(stake64, data)
    unknown.Scope.None()

    collection.Begin()
    unknown.Begin()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        collection.Add(key, uint64(index), key)
        unknown.Add(key, uint64(index), key)
    }

    collection.End()
    unknown.End()

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        proof, _ := collection.Get(key).Proof()
        if !(unknown.Verify(proof)) {
            test.Error("[verifiers.go]", "[verify]", "Verify() fails on valid proof.")
        }
    }

    ctx.verify.tree("[verify]", &unknown)

    for index := 0; index < 512; index++ {
        key := make([]byte, 8)
        binary.BigEndian.PutUint64(key, uint64(index))

        ctx.verify.values("[verify]", &unknown, key, uint64(index), key)
    }
}
