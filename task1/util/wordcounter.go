package util

import (
	"math"
	"strings"
)

// EngFreq is the frequency of each character in the English language
var EngFreq = map[rune]float64{
	'a': 0.08167, 'b': 0.01492, 'c': 0.02782, 'd': 0.04253, 'e': 0.12702, 'f': 0.02228, 'g': 0.02015, 'h': 0.06094,
	'i': 0.06966, 'j': 0.00153, 'k': 0.00772, 'l': 0.04025, 'm': 0.02406, 'n': 0.06749, 'o': 0.07507, 'p': 0.01929,
	'q': 0.00095, 'r': 0.05987, 's': 0.06327, 't': 0.09056, 'u': 0.02758, 'v': 0.00978, 'w': 0.02360, 'x': 0.00150,
	'y': 0.01974, 'z': 0.00074,
}

// WordCounter is a struct that stores the frequency of each character in a text
type WordCounter struct {
	Freq       map[rune]int
	text       string
	textLength int
}

// FreqMap is a map that stores the frequency of each character
type FreqMap map[rune]int

// NewWordCounter creates a new WordCounter
func NewWordCounter(text string, noFilter bool) *WordCounter {
	wc := &WordCounter{
		Freq: make(FreqMap),
		text: text,
	}

	if !noFilter {
		wc.text = strings.ToLower(wc.text)
		wc.text = FilterWords(wc.text)
	}

	wc.textLength = len(wc.text)

	return wc
}

// CountCharFreq counts the frequency of each character in the text
func (wc *WordCounter) CountCharFreq() {
	for _, char := range wc.text {
		wc.Freq[char]++
	}
}

// CalcConfidence calculates the confidence of the frequency of each character
func (wc *WordCounter) CalcConfidence() float64 {
	return CalcConfidence(wc.Freq, wc.textLength)
}

// CalcIC calculates the index of coincidence of the text
func CalcIC(freq FreqMap, textLength int) float64 {
	var ic float64
	for _, count := range freq {
		ic += float64(count) * float64(count-1)
	}
	ic /= float64(textLength) * float64(textLength-1)
	return ic
}

// CalcConfidence calculates the confidence of the frequency of each character
func CalcConfidence(freq FreqMap, textLength int) float64 {
	var confidence float64
	var appear int
	for c := 'a'; c <= 'z'; c++ {
		if engFreq, ok := EngFreq[c]; ok {
			freqChar := float64(freq[c]) / float64(textLength)
			confidence += math.Pow(engFreq-freqChar, 2)
			appear++
		}
	}
	if appear == 0 {
		return 0
	}
	return 1 / (1 + confidence/float64(appear))
}

// FinalConfidence calculates the final confidence of the text
func FinalConfidence(freq FreqMap, textLength int) float64 {
	engIC := 0.0667
	ic := CalcIC(freq, textLength)

	delta := math.Abs(ic - engIC)
	confidence := CalcConfidence(freq, textLength)
	icConfidence := 1 / (1 + delta)

	return (confidence + icConfidence) / 2
}
