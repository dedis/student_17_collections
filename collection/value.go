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

/*
    Data is a Value struct that can be used to store raw `[]byte` data along
    with each element of the collection (thus rendering it a key/value store).
    Data values are only stored in non-placeholder leaves, and don't propagate
    to parents. Therefore, Data values cannot be queried (an error will always
    be returned when calling `Navigate`).
*/
type Data struct {
}

// Interface

/*
    Placeholder returns an empty value (placeholder leaves don't need to store
    any raw data).
*/
func (this Data) Placeholder() []byte {
    return []byte{}
}

/*
    Parent returns an empty value (raw data is not propagated to interior
    nodes).
*/
func (this Data) Parent(left []byte, right []byte) []byte {
    return []byte{}
}

/*
    Navigate will always return an error (raw data cannot be navigated, as it
    is not propagated to interior nodes).
*/
func (this Data) Navigate(query []byte, parent []byte, left []byte, right []byte) (bool, error) {
    return false, errors.New("Data: data values cannot be navigated.")
}

// Stake64

/*
    Stake64 is a Value struct that can be used to associate to each element of
    a collection an amount of stake. Placeholder leaves have zero stake, and
    the stake of each internal node is given by the sum of the stakes of its
    children.

    This allows navigation by inversion of the cumulative stake: a query value
    uniformly distributed between 0 and the total stake on the tree (namely, the
    stake of the root) will yield an element of the collection with probability
    proportional to its stake.
*/
type Stake64 struct {
}

// Methods

/*
    Encode takes an `uint64` amount of stake and returns a `[]byte` value that
    can be stored on a node of the collection.
*/
func (this Stake64) Encode(stake uint64) []byte {
    value := make([]byte, 8)
    binary.BigEndian.PutUint64(value, stake)
    return value
}

/*
    Decode takes a `[]byte` value (e.g., stored by a node of the collection) and
    returns the `uint64` amount of stake that it encodes.
*/
func (this Stake64) Decode(value []byte) uint64 {
    return binary.BigEndian.Uint64(value)
}

// Interface

/*
    Placeholder returns an encoded zero-stake value (placeholder leaves always
    have zero stake).
*/
func (this Stake64) Placeholder() []byte {
    return this.Encode(0)
}

/*
    Parent returns a value encoding a stake equal to the sum of the stakes
    encoded by the children leaves (the stake on each node represents the total
    amount of stake under it; this allows cumulative inversion navigation).
*/
func (this Stake64) Parent(left []byte, right []byte) []byte {
    return this.Encode(this.Decode(left) + this.Decode(right))
}

/*
    Navigate implements cumulative inversion stake navigation: if the query
    stake is greater or equal to the parent stake, it yields an error.
    Otherwise, it navigates left (returns `false`) if the query stake is smaller
    than the left stake, and right (returns `true`) otherwise. When navigating
    right, the query stake is decreased by the left stake.
*/
func (this Stake64) Navigate(query []byte, parent []byte, left []byte, right []byte) (bool, error) {
    querystake := this.Decode(query)
    parentstake := this.Decode(parent)

    if querystake >= parentstake {
        return false, errors.New("Stake64: query stake exceeds parent stake.")
    }

    leftstake := this.Decode(left)

    if querystake >= leftstake {
        copy(query, this.Encode(querystake - leftstake))
        return true, nil
    } else {
        return false, nil
    }
}
