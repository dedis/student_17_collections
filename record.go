package collection

import "errors"

type Record struct {
    collection *collection
    match bool

    key []byte
    values [][]byte
}

// Constructors

func recordmatch(collection *collection, node *node) Record {
    return Record{collection, true, node.key, node.values}
}

func recordmismatch(collection *collection, key []byte) Record {
    return Record{collection, false, key, [][]byte{}}
}

// Getters

func (this *Record) Match() bool {
    return this.match
}

func (this *Record) Key() []byte {
    return this.key
}

func (this *Record) Values() ([]interface{}, error) {
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
