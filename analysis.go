package gotar

import (
	"math"
	"time"
)

type note struct {
	name      string
	frequency float64
}

var noteMap = []string{"A", "B♭", "B", "C", "C#", "D", "E♭", "E", "F", "F#", "G", "G#"}

// FrequencyInfo information about a frequency
type FrequencyInfo struct {
	Frequency float64
	Note      string
	Octave    int
}

var unknownFrequency = FrequencyInfo{0.0, "", 0}

var notes = []note{
	note{"E", 82.4069},
	note{"A", 110},
	note{"D", 146.832},
	note{"G", 195.998},
	note{"B", 246.932},
	note{"E", 329.628},
}

var assessStringsUntilTime int64 = 0
var assessedStringsInLastFrame = false

const rmsThreshold = 0.006

var lastRms = 0.0

var differences = make([]float64, len(notes))
var lastMinDifference = 0

// GetFrequencyInfo get frequency information about a waveform
func GetFrequencyInfo(waveform []float64, sampleRate float64) FrequencyInfo {

	searchSize := len(waveform) / 2

	tolerance := 0.001
	rms := 0.0
	rmsMin := 0.008

	prevAssessedStrings := assessedStringsInLastFrame

	for _, amplitude := range waveform {
		rms += amplitude * amplitude
	}

	rms = math.Sqrt(rms / float64(len(waveform)))

	if rms < rmsMin {
		return unknownFrequency
	}

	time := (time.Now().UnixNano() / 1000000)

	if rms > lastRms+rmsThreshold {
		assessStringsUntilTime = time + 250
	}

	if time < assessStringsUntilTime {
		assessedStringsInLastFrame = true

		for i, note := range notes {
			offset := int(math.Round(sampleRate / note.frequency))
			difference := 0.0

			if !prevAssessedStrings {
				differences[i] = 0
			}

			for j := 0; j < searchSize; j++ {
				currentAmp := waveform[j]
				offsetAmp := waveform[j+offset]
				difference += math.Abs(currentAmp - offsetAmp)
			}

			difference /= float64(searchSize)

			differences[i] += difference * float64(offset)
		}
	} else {
		assessedStringsInLastFrame = false
	}

	if !assessedStringsInLastFrame && prevAssessedStrings {
		lastMinDifference = argmin(differences)
	}

	assumedString := notes[lastMinDifference]
	searchRange := 10
	actualFrequency := int(math.Round(sampleRate / assumedString.frequency))
	searchStart := actualFrequency - searchRange
	searchEnd := actualFrequency + searchRange
	smallestDifference := math.Inf(1)

	for i := searchStart; i < searchEnd; i++ {
		difference := 0.0

		for j := 0; j < searchSize; j++ {
			currentAmp := waveform[j]
			offsetAmp := waveform[j+i]
			difference += math.Abs(currentAmp - offsetAmp)
		}

		difference /= float64(searchSize)

		if difference < smallestDifference {
			smallestDifference = difference
			actualFrequency = i
		}

		if difference < tolerance {
			actualFrequency = i
			break
		}

	}

	lastRms = rms

	frequency := sampleRate / float64(actualFrequency)
	octave := getOctave(frequency)
	note := getNote(frequency)

	return FrequencyInfo{frequency, note, octave}
}

func linearizeFreq(frequency float64) float64 {
	return math.Log2(frequency / 440.0)
}

func getOctave(frequency float64) int {
	semitonesFromA4 := 12 * linearizeFreq(frequency)
	octave := 4 + ((9 + semitonesFromA4) / 12)
	octave = math.Floor(octave)
	return int(octave)
}

func getNote(frequency float64) string {
	semitonesFromA4 := 12 * linearizeFreq(frequency)
	noteNum := (12 + (int(math.Round(semitonesFromA4)) % 12)) % 12
	return noteMap[noteNum]
}

func argmin(arr []float64) int {
	min := 0

	for i, value := range arr {
		if value < arr[min] {
			min = i
		}
	}

	return min
}
