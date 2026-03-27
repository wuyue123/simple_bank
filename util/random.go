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
	"strings"
	"time"
)

var alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandomString generates a random string of the given length.
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomInt generates a random integer between 0 and n-1.
func RandomInt(n int) int {
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	return rd.Intn(n)
}

// RandomMoney generates a random amount of money between 0 and 1000.
func RandomMoney() int64 {
	return int64(RandomInt(1000))
}

func RandomEmail() string {
	return RandomString(6) + "@163.com"
}

// RandomCurrency returns a random currency code.
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	return currencies[RandomInt(len(currencies))]
}
