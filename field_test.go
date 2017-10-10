package collection

import "testing"
import "math/rand"

func TestFieldData(test *testing.T) {
    var data Data

    if !equal(data.Encode([]byte("mydata")), []byte("mydata")) {
        test.Error("[field.go]", "[encode]", "Data encode() function is not the identity function.")
    }

    if !equal(data.Decode([]byte("mydata")).([]byte), []byte("mydata")) {
        test.Error("[field.go]", "[encode]", "Data decode() function is not the identity function.")
    }

    if len(data.Placeholder()) != 0 {
        test.Error("[field.go]", "[placeholder]", "Non-empty data placeholder.")
    }

    if len(data.Parent([]byte("leftdata"), []byte("rightdata"))) != 0 {
        test.Error("[field.go]", "[parent]", "Non-empty data parent: Data should not propagate anything to its parent.")
    }

    _, err := data.Navigate([]byte("query"), []byte("parentdata"), []byte("leftdata"), []byte("rightdata"))
    if err == nil {
        test.Error("[field.go]", "[navigate]", "Data navigation does not yield errors. It should be impossible to navigate Data.")
    }
}

func TestFieldStake64(test *testing.T) {
    var stake64 Stake64

    for trial := 0; trial < 64; trial++ {
        stake := rand.Uint64()
        if stake != stake64.Decode(stake64.Encode(stake)).(uint64) {
            test.Error("[field.go]", "[encodeconsistency]", "Stake64 ncode / decode inconsistency.")
        }
    }

    if stake64.Decode(stake64.Placeholder()).(uint64) != 0 {
        test.Error("[field.go]", "[placeholder]", "Non-zero placeholder stake.")
    }

    for trial := 0; trial < 64; trial++ {
        leftstake := rand.Uint64()
        rightstake := rand.Uint64()

        left := stake64.Encode(leftstake)
        right := stake64.Encode(rightstake)

        if stake64.Decode(stake64.Parent(left, right)).(uint64) != (leftstake + rightstake) {
            test.Error("[field.go]", "[parent]", "Parent stake is not equal to the sum of children stakes.")
        }
    }

    for trial := 0; trial < 64; trial++ {
        leftstake := uint64(rand.Uint32())
        rightstake := uint64(rand.Uint32())
        parentstake := leftstake + rightstake

        parent := stake64.Encode(parentstake)
        left := stake64.Encode(leftstake)
        right := stake64.Encode(rightstake)

        wrongquerystake := parentstake + uint64(rand.Uint32())
        wrongquery := stake64.Encode(wrongquerystake)

        navigation, err := stake64.Navigate(wrongquery, parent, left, right)

        if err == nil {
            test.Error("[field.go]", "[navigate]", "Error not yielded on illegal Stake64 query.")
        }

        if stake64.Decode(wrongquery).(uint64) != wrongquerystake {
            test.Error("[field.go]", "[navigate]", "Navigate altered illegal Stake64 query without navigating.")
        }

        querystake := rand.Uint64() % parentstake
        query := stake64.Encode(querystake)

        navigation, err = stake64.Navigate(query, parent, left, right)

        if err != nil {
            test.Error("[field.go]", "[navigate]", "Error yielded on legal Stake64 query.")
        }

        newquerystake := stake64.Decode(query).(uint64)

        if querystake >= leftstake {
            if !navigation {
                test.Error("[field.go]", "[navigate]", "Stake64 navigates on wrong child.")
            }

            if newquerystake != querystake - leftstake {
                test.Error("[field.go]", "[navigate]", "Stake64 right navigation doesn't correctly decrease the query.")
            }
        } else {
            if navigation {
                test.Error("[field.go]", "[navigate]", "Stake64 navigates on wrong child.")
            }

            if newquerystake != querystake {
                test.Error("[field.go]", "[navigate]", "Stake64 left navigation altered the query.")
            }
        }
    }
}
