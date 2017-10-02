package collection

import "testing"
import "math/rand"

func TestStake64(test *testing.T) {
    var stake64 Stake64

    for trial := 0; trial < 64; trial++ {
        stake := rand.Uint64()
        if stake != stake64.Decode(stake64.Encode(stake)) {
            test.Error("[encodeconsistency]", "Encode / decode inconsistency.")
        }
    }

    if stake64.Decode(stake64.Placeholder()) != 0 {
        test.Error("[placeholder]", "Non-zero placeholder stake.")
    }

    for trial := 0; trial < 64; trial++ {
        leftstake := rand.Uint64()
        rightstake := rand.Uint64()

        left := stake64.Encode(leftstake)
        right := stake64.Encode(rightstake)

        if stake64.Decode(stake64.Parent(left, right)) != (leftstake + rightstake) {
            test.Error("[parent]", "Parent stake is not equal to the sum of children stakes.")
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
            test.Error("[navigate]", "Error not yielded on illegal query.")
        }

        if stake64.Decode(wrongquery) != wrongquerystake {
            test.Error("[navigate]", "Navigate altered illegal query without navigating.")
        }

        querystake := rand.Uint64() % parentstake
        query := stake64.Encode(querystake)

        navigation, err = stake64.Navigate(query, parent, left, right)

        if err != nil {
            test.Error("[navigate]", "Error yielded on legal query.")
        }

        newquerystake := stake64.Decode(query)

        if querystake >= leftstake {
            if !navigation {
                test.Error("[navigate]", "Navigated on wrong child.")
            }

            if newquerystake != querystake - leftstake {
                test.Error("[navigate]", "Right navigation doesn't correctly decrease the query.")
            }
        } else {
            if navigation {
                test.Error("[navigate]", "Navigated on wrong child.")
            }

            if newquerystake != querystake {
                test.Error("[navigate]", "Left navigation altered the query.")
            }
        }
    }
}
