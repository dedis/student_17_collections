package collection

import "errors"

type Record struct {
    collection *collection

    field int
    query []byte
    match bool

    key []byte
    values [][]byte
}

// Constructors

func recordkeymatch(collection *collection, node *node) Record {
    return Record{collection, 0, []byte{}, true, node.key, node.values}
}

func recordquerymatch(collection *collection, field int, query interface{}, node *node) Record {
    if field >= len(collection.fields) {
        panic("Field out of range.")
    }

    return Record{collection, field, collection.fields[field].Encode(query), true, node.key, node.values}
}

func recordkeymismatch(collection *collection, key []byte) Record {
    return Record{collection, 0, []byte{}, false, key, [][]byte{}}
}

// Getters

func (this Record) Query() (interface{}, error) {
    if len(this.query) == 0 {
        return nil, errors.New("No query specified.")
    }

    if len(this.values) <= this.field {
        return nil, errors.New("Field out of range.")
    }

    value, err := this.collection.fields[this.field].Decode(this.query)

    if err != nil {
        return nil, err
    }

    return value, nil
}

func (this Record) Match() bool {
    return this.match
}

func (this Record) Key() []byte {
    return this.key
}

func (this Record) Values() ([]interface{}, error) {
    if !(this.match) {
        return []interface{}{}, errors.New("No match found.")
    }

    if len(this.values) != len(this.collection.fields) {
        return []interface{}{}, errors.New("Wrong number of values.")
    }

    var values []interface{}

    for index := 0; index < len(this.values); index++ {
        value, err := this.collection.fields[index].Decode(this.values[index])

        if err != nil {
            return []interface{}{}, err
        }

        values = append(values, value)
    }

    return values, nil
}
