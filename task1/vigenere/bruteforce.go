package vigenere

import (
	"USST-CyberSecurity/task1/kasiski"
	"USST-CyberSecurity/task1/util"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type FreqMap = util.FreqMap

// PossibleKey represents a possible key for the Vigenere cipher
type PossibleKey struct {
	possibility float64
	key         string
	freq        FreqMap
}

// Key represents a key for the Vigenere cipher
type Key struct {
	Possibility float64
	Key         string
}

// BruteForceDecrypt decrypts the given cipher with known key length
func BruteForceDecrypt(cipher string, keyLength int) []Key {
	// Possible key of every position has a list
	keyPossibilities := make([][]PossibleKey, keyLength)

	var dfs func(result *[]Key, l int, curr PossibleKey, s FreqMap)

	dfs = func(result *[]Key, l int, curr PossibleKey, ori FreqMap) {
		s := make(FreqMap) // make sure to copy the map
		// merge the frequency map of the current key
		for c := 'a'; c <= 'z'; c++ {
			s[c] = ori[c] + curr.freq[c]
		}

		if l == keyLength {
			// Calculate both Ic and confidence to get the final confidence
			conf := util.FinalConfidence(s, len(cipher))
			*result = append(*result, Key{conf, curr.key})
			return
		}

		for _, i := range keyPossibilities[l] {
			// make a copy of the key
			i1 := PossibleKey{0, "", make(FreqMap)}
			for c := 'a'; c <= 'z'; c++ {
				i1.freq[c] = i.freq[c]
			}
			i1.possibility = i.possibility
			i1.key = curr.key + i.key

			dfs(result, l+1, i1, s)
		}
	}

	cipher = util.FilterWords(strings.ToLower(cipher))
	slices := make([]string, keyLength)
	for i := range slices {
		for j := i; j < len(cipher); j += keyLength {
			slices[i] += string(cipher[j])
		}
	}

	// For each slice, try all the possible keys
	for i, s := range slices {
		possibilities := make([]PossibleKey, 0, 26)
		for shift := 0; shift < 26; shift++ {
			key := string('a' + rune(shift))
			decrypted := Decrypt(s, key)
			wc := util.NewWordCounter(decrypted, true)
			wc.CountCharFreq()
			possibility := wc.CalcConfidence()
			possibilities = append(possibilities, PossibleKey{possibility, key, wc.Freq})
		}

		// Sort by confidence, descending
		sort.Slice(possibilities, func(i, j int) bool {
			return possibilities[i].possibility > possibilities[j].possibility
		})

		// Only keep the top 3 possibilities
		if len(possibilities) > 3 {
			possibilities = possibilities[:3]
		}

		keyPossibilities[i] = possibilities
	}

	results := make([]Key, 0)
	dfs(&results, 0, PossibleKey{0, "", make(FreqMap)}, make(FreqMap))

	// Sort by confidence, descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Possibility > results[j].Possibility
	})

	if len(results) > 5 {
		results = results[:5]
	}

	return results
}

// MultiThreadDecrypt decrypts the given cipher without key
func MultiThreadDecrypt(cipher string) []Key {
	// do the Kasiski test to find the possible key length
	cipherSliceLen := min(1000, len(cipher))
	possibleLengthList := kasiski.Test(cipher[:cipherSliceLen], 3)

	// only keep the top 5 possible key lengths
	possibleLengthListSliceLen := min(5, len(possibleLengthList))
	possibleLengthList = possibleLengthList[:possibleLengthListSliceLen]

	// if no possible key length is found, try all possible key lengths
	if len(possibleLengthList) == 0 {
		for i := 2; i <= min(12, len(cipher)); i++ { // only try key length up to 12
			possibleLengthList = append(possibleLengthList, i)
		}
	}

	fmt.Printf("Possible key length list: %v\n", possibleLengthList)

	var totalResult []Key

	writeLock := sync.Mutex{}
	wg := sync.WaitGroup{}

	// try all possible key lengths
	for _, length := range possibleLengthList {
		if length > 14 {
			continue
		}
		wg.Add(1)
		go func(length int) {
			results := BruteForceDecrypt(cipher, length)
			writeLock.Lock()
			for _, result := range results {
				totalResult = append(totalResult, result)
			}
			writeLock.Unlock()
			fmt.Printf("Finished key length: %d\n", length)
			wg.Done()
		}(length)
	}

	wg.Wait()

	// sort by confidence, descending
	sort.Slice(totalResult, func(i, j int) bool {
		return totalResult[i].Possibility > totalResult[j].Possibility
	})

	return totalResult
}
