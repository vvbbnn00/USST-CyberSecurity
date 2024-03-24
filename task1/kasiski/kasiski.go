package kasiski

import (
	"USST-CyberSecurity/task1/util"
	"sort"
	"strings"
)

// FindRepeatedSequences finds repeated sequences in the cipher
func FindRepeatedSequences(cipher string, minSeqLength int) map[string][]int {
	sequences := make(map[string][]int)

	// Count the number of repeated sequences
	for i := 0; i <= len(cipher)-minSeqLength; i++ {
		seq := cipher[i : i+minSeqLength]
		sequences[seq] = append(sequences[seq], i)
	}

	// Remove sequences that are not repeated
	for seq, positions := range sequences {
		if len(positions) < 2 {
			delete(sequences, seq)
		}
	}

	return sequences
}

// Test performs the Kasiski test on the cipher
func Test(cipher string, minSeqLength int) []int {
	cipher = util.FilterWords(strings.ToLower(cipher))

	repeatedSeqs := FindRepeatedSequences(cipher, minSeqLength)
	var distances []int

	// Calculate the distances between the repeated sequences
	for _, positions := range repeatedSeqs {
		for i := 1; i < len(positions); i++ {
			distances = append(distances, positions[i]-positions[i-1])
		}
	}

	// Calculate the factors of the distances
	factors := make(map[int]int)
	for _, distance := range distances {
		for i := 2; i <= distance; i++ {
			if distance%i == 0 {
				factors[i]++
			}
		}
	}

	// Sort the factors by the number of occurrences
	sortedFactors := make([]int, 0, len(factors))
	for factor := range factors {
		sortedFactors = append(sortedFactors, factor)
	}
	sort.Slice(sortedFactors, func(i, j int) bool {
		return factors[sortedFactors[i]] > factors[sortedFactors[j]]
	})

	return sortedFactors
}
