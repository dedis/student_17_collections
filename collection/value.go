package collection

import "errors"
import "encoding/binary"

// Interfaces

type value interface {
    Placeholder() []byte
    Parent([]byte, []byte) []byte
    Navigate([]byte, []byte, []byte, []byte) (bool, error)
}

// Structs

// Data

type Data struct {
}

// Interface

func (this Data) Placeholder() []byte {
    return []byte{}
}

func (this Data) Parent(left []byte, right []byte) []byte {
    return []byte{}
}

func (this Data) Navigate(query []byte, parent []byte, left []byte, right []byte) (bool, error) {
    return false, errors.New("Data: data values cannot be navigated.")
}

// Stake64

type Stake64 struct {
}

// Methods

func (this Stake64) Encode(stake uint64) []byte {
    value := make([]byte, 8)
    binary.BigEndian.PutUint64(value, stake)
    return value
}

func (this Stake64) Decode(value []byte) uint64 {
    return binary.BigEndian.Uint64(value)
}

// Interface

func (this Stake64) Placeholder() []byte {
    return make([]byte, 8)
}

func (this Stake64) Parent(left []byte, right []byte) []byte {
    leftstake := binary.BigEndian.Uint64(left)
    rightstake := binary.BigEndian.Uint64(right)

    parentstake := leftstake + rightstake
    parent := make([]byte, 8)
    binary.BigEndian.PutUint64(parent, parentstake)

    return parent
}

func (this Stake64) Navigate(query []byte, parent []byte, left []byte, right []byte) (bool, error) {
    querystake := binary.BigEndian.Uint64(query)
    parentstake := binary.BigEndian.Uint64(parent)

    if querystake >= parentstake {
        return false, errors.New("Stake64: query stake exceeds parent stake.")
    }

    leftstake := binary.BigEndian.Uint64(left)

    if querystake >= leftstake {
        binary.BigEndian.PutUint64(query, querystake - leftstake)
        return true, nil
    } else {
        return false, nil
    }
}
