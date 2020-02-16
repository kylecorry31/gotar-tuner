package gotar

import "math"

var noteMap = []string{"A", "B♭", "B", "C", "C#", "D", "E♭", "E", "F", "F#", "G", "G#"}

// FrequencyInfo information about a frequency
type FrequencyInfo struct {
	Frequency float64
	Note      string
	Octave    int
}

// UnknownFrequency an unknown frequency
var UnknownFrequency = FrequencyInfo{0.0, "", 0}

func linearizeFrequency(frequency float64) float64 {
	return math.Log2(frequency / 440.0)
}

// CreateFrequencyInfo create the frequency information
func CreateFrequencyInfo(frequency float64) FrequencyInfo {
	return FrequencyInfo{frequency, GetNote(frequency), GetOctave(frequency)}
}

// GetOctave the octave from a frequency
func GetOctave(frequency float64) int {
	semitonesFromA4 := 12 * linearizeFrequency(frequency)
	octave := 4 + ((9 + semitonesFromA4) / 12)
	note := GetNote(frequency)
	if note == "C" {
		octave = math.Round(octave)
	} else {
		octave = math.Floor(octave)
	}
	return int(octave)
}

// GetNote get the note from a frequency
func GetNote(frequency float64) string {
	semitonesFromA4 := 12 * linearizeFrequency(frequency)
	noteNum := (12 + (int(math.Round(semitonesFromA4)) % 12)) % 12
	return noteMap[noteNum]
}
