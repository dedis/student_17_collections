package collection

import "testing"

func TestUpdateProxy(test *testing.T) {
    collection := EmptyCollection()

    proxy := collection.proxy([][]byte{[]byte("firstkey"), []byte("secondkey"), []byte("thirdkey")})

    if proxy.collection != &collection {
        test.Error("[update.go]", "[proxy]", "proxy() method sets wrong collection pointer.")
    }

    if !(proxy.paths[sha256([]byte("firstkey"))]) || !(proxy.paths[sha256([]byte("secondkey"))]) || !(proxy.paths[sha256([]byte("thirdkey"))]) {
        test.Error("[update.go]", "[proxy]", "proxy() method does not set the paths provided.")
    }

    if proxy.paths[sha256([]byte("otherkey"))] {
        test.Error("[update.go]", "[proxy]", "proxy() sets more paths than the ones provided.")
    }
}

func TestUpdateProxyMethods(test *testing.T) {
    ctx := testctx("[update.go]", test)

    stake64 := Stake64{}
    collection := EmptyCollection(stake64)

    proxy := collection.proxy([][]byte{[]byte("firstkey"), []byte("secondkey"), []byte("thirdkey")})

    collection.Add([]byte("firstkey"), uint64(66))
    record := proxy.Get([]byte("firstkey"))

    if !(record.Match()) {
        test.Error("[update.go]", "[get]", "Proxy method get() does not return an existing record.")
    }

    values, _ := record.Values()
    if values[0].(uint64) != 66 {
        test.Error("[update.go]", "[get]", "Proxy method get() returns wrong values.")
    }

    error := proxy.Add([]byte("secondkey"), uint64(33))
    if error != nil {
        test.Error("[update.go]", "[add]", "Proxy method add() yields an error on valid key.")
    }

    error = proxy.Add([]byte("secondkey"), uint64(22))
    if error == nil {
        test.Error("[update.go]", "[add]", "Proxy method add() does not yield an error when adding an existing key.")
    }

    record, _ = collection.Get([]byte("secondkey")).Record()

    if !(record.Match()) {
        test.Error("[update.go]", "[add]" ,"Proxy method add() does not add the record provided.")
    }

    values, _ = record.Values()
    if values[0].(uint64) != 33 {
        test.Error("[update.go]", "[add]", "Proxy method add() adds wrong values.")
    }

    error = proxy.Set([]byte("secondkey"), uint64(22))
    if error != nil {
        test.Error("[update.go]", "[set]", "Proxy method set() yields an error when setting on an existing key.")
    }

    record, _ = collection.Get([]byte("secondkey")).Record()
    values, _ = record.Values()

    if values[0].(uint64) != 22 {
        test.Error("[update.go]", "[set]", "Proxy method set() does not set the correct values.")
    }

    error = proxy.Set([]byte("thirdkey"), uint64(11))
    if error == nil {
        test.Error("[update.go]", "[set]", "Proxy method set() does not yield an error when setting on a non-existing key.")
    }

    error = proxy.SetField([]byte("firstkey"), 0, uint64(11))
    if error != nil {
        test.Error("[update.go]", "[setfield]", "Proxy method setfield() does yields an error when setting on an existing key.")
    }

    record, _ = collection.Get([]byte("firstkey")).Record()
    values, _ = record.Values()

    if values[0].(uint64) != 11 {
        test.Error("[update.go]", "[setfield]", "Proxy method setfield() does not set the correct values.")
    }

    error = proxy.SetField([]byte("thirdkey"), 0, uint64(99))
    if error == nil {
        test.Error("[update.go]", "[setfield]", "Proxy method setfield() does not yield an error when setting on a non-existing key.")
    }

    error = proxy.Remove([]byte("secondkey"))
    if error != nil {
        test.Error("[update.go]", "[remove]", "Proxy method remove() yields an error when removing an existing key.")
    }

    record, _ = collection.Get([]byte("secondkey")).Record()
    if record.Match() {
        test.Error("[update.go]", "[remove]", "Proxy method remove() does not remove an existing key.")
    }

    error = proxy.Remove([]byte("secondkey"))
    if error == nil {
        test.Error("[update.go]", "[remove]", "Proxy method remove() does not yield an error when removing a non-existing key.")
    }

    ctx.should_panic("[get]", func() {
        proxy.Get([]byte("otherkey"))
    })

    ctx.should_panic("[add]", func() {
        proxy.Add([]byte("otherkey"), uint64(12))
    })

    ctx.should_panic("[set]", func() {
        proxy.Set([]byte("otherkey"), uint64(12))
    })

    ctx.should_panic("[setfield]", func() {
        proxy.SetField([]byte("otherkey"), 0, uint64(12))
    })

    ctx.should_panic("[remove]", func() {
        proxy.Remove([]byte("otherkey"))
    })
}

func TestUpdateProxyHas(test *testing.T) {
    collection := EmptyCollection()
    proxy := collection.proxy([][]byte{[]byte("firstkey"), []byte("secondkey"), []byte("thirdkey")})

    if !(proxy.has([]byte("firstkey"))) || !(proxy.has([]byte("secondkey"))) || !(proxy.has([]byte("thirdkey"))) {
        test.Error("[update.go]", "[has]", "Proxy method has() returns false on whitelisted key.")
    }

    if proxy.has([]byte("otherkey")) {
        test.Error("[update.go]", "[has]", "Proxy method has8) returns true on non-whitelisted key.")
    }
}
