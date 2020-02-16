package gotar

import (
	"math"
	"time"
)

type note struct {
	name      string
	frequency float64
}

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

// AutocorrelateFrequency get frequency information about a waveform
func AutocorrelateFrequency(waveform []float64, sampleRate float64) float64 {

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
		return 0
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

	return sampleRate / float64(actualFrequency)
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
