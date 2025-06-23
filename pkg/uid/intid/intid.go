package intid

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

var seed = time.Now().UnixNano()
var rnd = rand.New(rand.NewSource(seed))

type object struct {
	len int
}

func NewIntIDGenerator() *object {
	return &object{
		len: 20,
	}
}

func (o *object) SetIDLen(l int) {
	o.len = l
}

func (o *object) GenUID() string {
	return GenIDStr(o.len)
}

func GenInt(min int64, max int64) int64 {

	return rnd.Int63n(max-min) + min
}

func GenIDWithLen(minLen, maxLen uint) int64 {
	//int64 max +- 9223372036854775807
	if minLen >= maxLen {
		maxLen = minLen
	}
	if maxLen > 17 {
		maxLen = 17
	}
	if minLen > 17 {
		minLen = 17
	}

	if minLen >= maxLen {
		maxLen = minLen
	}

	var len int32
	if maxLen == minLen {
		len = int32(maxLen)
	} else {
		len = rnd.Int31n(int32(maxLen-minLen)) + int32(minLen)
	}

	var result int64
	for i := len; i >= 0; i-- {
		r1 := rnd.Int63n(9)
		if i == len {
			for r1 == 0 {
				r1 = rnd.Int63n(9)
			}
		}
		result += int64(r1 * int64(math.Pow(10, float64(i))))
	}

	return result
}

func GenIDStr(length int) string {
	var nid string

	for i := 0; i < length; i++ {
		var (
			r1 int
		)

		r1 = rnd.Intn(9)

		if i == 0 {
			for r1 == 0 {
				r1 = rnd.Intn(9)
			}
		}

		nid += strconv.Itoa(r1)
	}
	return nid
}

func SetSeed(s int64) {
	seed = s
	rnd = rand.New(rand.NewSource(seed))
}
