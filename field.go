package collection

import "errors"
import "encoding/binary"

// Enums

type Navigation bool

const(
    Left Navigation = false
    Right Navigation = true
)

// Interfaces

type Field interface {
    Encode(interface{}) []byte
    Decode([]byte) interface{}
    Placeholder() []byte
    Parent([]byte, []byte) []byte
    Navigate([]byte, []byte, []byte, []byte) (Navigation, error)
}

// Structs

// Data

type Data struct {
}

// Interface

func (this Data) Encode(generic interface{}) []byte {
    value := generic.([]byte)
    return value
}

func (this Data) Decode(raw []byte) interface{} {
    return raw
}

func (this Data) Placeholder() []byte {
    return []byte{}
}

func (this Data) Parent(left []byte, right []byte) []byte {
    return []byte{}
}

func (this Data) Navigate(query []byte, parent []byte, left []byte, right []byte) (Navigation, error) {
    return false, errors.New("Data values cannot be navigated.")
}

// Stake64

type Stake64 struct {
}

// Interface

func (this Stake64) Encode(generic interface{}) []byte {
    value := generic.(uint64)
    raw := make([]byte, 8)

    binary.BigEndian.PutUint64(raw, value)
    return raw
}

func (this Stake64) Decode(raw []byte) interface{} {
    return binary.BigEndian.Uint64(raw)
}

func (this Stake64) Placeholder() []byte {
    return this.Encode(uint64(0))
}

func (this Stake64) Parent(left []byte, right []byte) []byte {
    return this.Encode(this.Decode(left).(uint64) + this.Decode(right).(uint64))
}

func (this Stake64) Navigate(query []byte, parent []byte, left []byte, right []byte) (Navigation, error) {
    queryvalue := this.Decode(query).(uint64)
    parentvalue := this.Decode(parent).(uint64)

    if queryvalue >= parentvalue {
        return false, errors.New("Query exceeds parent stake.")
    }

    leftvalue := this.Decode(left).(uint64)

    if queryvalue >= leftvalue {
        copy(query, this.Encode(queryvalue - leftvalue))
        return Right, nil
    } else {
        return Left, nil
    }
}
