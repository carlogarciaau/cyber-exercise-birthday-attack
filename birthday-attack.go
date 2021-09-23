package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		println("Usage: ./birthday-attack confession_fake.txt confession_real.txt")
		os.Exit(0)
	}

	fake, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fakeStr := string(fake)

	real, err2 := ioutil.ReadFile(os.Args[2])
	if err2 != nil {
		log.Fatal(err2)
	}

	realStr := string(real)
	var fakeList []string
	var realList []string

	for i := 1; i <= 1000; i++ {
		for j := 1; j <= 23; j++ {
			newFakeStr := replaceNth(fakeStr, "\n", addSpaces(i), j)
			fakeList = append(fakeList, newFakeStr)

			newRealStr := replaceNth(realStr, "\n", addSpaces(i), j)
			realList = append(realList, newRealStr)
		}
	}

	for i := 0; i < len(fakeList); i++ {
		sumFake := sha256.Sum256([]byte(fakeList[i]))
		sumFakeStr := hex.EncodeToString(sumFake[:])

		for j := 0; j < len(realList); j++ {
			sumReal := sha256.Sum256([]byte(realList[j]))
			sumRealStr := hex.EncodeToString(sumReal[:])

			if findCollision(sumFakeStr, sumRealStr) != "" {
				println(sumFakeStr + " - " + sumRealStr)
				writeMessageToFile(fakeList[i], "confession_fake.txt."+sumFakeStr)
				writeMessageToFile(realList[j], "confession_real.txt."+sumRealStr)
			}
		}
	}
}

// Find collision at the end of each hash that is X or more hex digits
func findCollision(hash1, hash2 string) string {
	hashMatch := ""
	hexDigitCollisionCriteria := 8

	for i := hexDigitCollisionCriteria; i <= len(hash1); i++ {
		if hash1[len(hash1)-i:len(hash1)] == hash2[len(hash2)-i:len(hash2)] {
			hashMatch = hash1[len(hash1)-i : len(hash1)]
			fmt.Println("Collision Criteria: " + strconv.Itoa(hexDigitCollisionCriteria) + " hex digits.")
		} else {
			break
		}
	}

	return hashMatch
}

// Replace the nth occurrence of old in s by new.
func replaceNth(s, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return s[:i] + new + s[i+len(old):]
		}
		i += len(old)
	}
	return s
}

// Add any number of spaces
func addSpaces(length int) string {
	returnValue := ""
	for i := 1; i <= length; i++ {
		returnValue = returnValue + " "
	}

	return returnValue + "\n"
}

// Output altered messages to a file
func writeMessageToFile(message, filename string) {
	data := []byte(message)
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
