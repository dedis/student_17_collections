package collection

import "errors"
import "encoding/binary"

// Interfaces

/*
    Value is the interface that is used to define classes of values that can be
    stored on the nodes of a collection.

    A Value struct should offer the following methods:
     - `Placeholder() []byte`, that returns the value that should be stored on
       a placeholder leaf.
     - `Parent(left []byte, right []byte) []byte`, that returns the value that
       should be stored on an internal node, given the value of its left and
       right children.
     - `Navigate(query []byte, parent []byte, left []byte, right []byte) (bool, error)`,
       that, given a query value, the value of an internal node and that of its
       left and right children, returns `(true, nil)` if the query can be
       satisfied by navigating right, `(false, nil)` if the query can be
       satisfied by navigating left, and `(?, err)` if the query cannot be
       satisfied.

    For example, two Value structs that collection offers are:
     - `Data`, that can be used to store raw `[]byte` data with each element of
       the collection (which will therefore serve also as a key/value store).
       In this case, both `Placeholder` and `Parent` will return an empty value
       (placeholder leaves should carry no data and data should not be
       propagated to internal nodes). `Navigate` will always return an error
       (as raw data is not designed to be navigated).
     - `Stake64`, that can be used to store a quantity of stake owned by each
       identifier in the collection. In this case, the value of a placeholder
       leaf is always zero, and the value of an internal node is equal to the
       sum of the values of its left and right children. `Navigate` will
       accept a query between 0 and the value of the parent leaf, and navigate
       left if the query is smaller than the value of the left child, and
       right otherwise. When navigating right, the query value will be decreased
       by the value of the left leaf. Recursively navigating a `Stake64` from
       the root with a random query uniformly distributed between 0 and the
       value of the root (i.e., the total stake in the collection) will produce
       a random node, with probability proportional to its stake.
*/
type Value interface {
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
