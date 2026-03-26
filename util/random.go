/**
* @description:
* @author
* @date 2026-03-24 23:50:27
* @version 1.0
*
* Change Logs:
* Date           Author       Notes
*
 */

package util

import (
	"math/rand"
	"time"
)


var alphaNumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandomString generates a random string of the given length.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphaNumeric[RandomInt(len(alphaNumeric))]
	}
	return string(b)
}

// RandomInt generates a random integer between 0 and n-1.
func RandomInt(n int) int {
	rd:=rand.New(rand.NewSource(time.Now().Unix()))
	return rd.Intn(n) 
}

// RandomMoney generates a random amount of money between 0 and 1000.
func RandomMoney() int64 {
	return int64(RandomInt(1000))
}

// RandomCurrency returns a random currency code.
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[RandomInt(len(currencies))]
}