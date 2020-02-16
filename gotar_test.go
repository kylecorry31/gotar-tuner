package gotar

import (
	"math"
	"testing"
)

func TestCreateFrequencyInfo(t *testing.T) {
	tables := []struct {
		frequency float64
		note      string
		octave    int
	}{
		{16.35, "C", 0},
		{32.7, "C", 1},
		{65.41, "C", 2},
		{2093, "C", 7},
		{69.3, "C#", 2},
		{146.8, "D", 3},
		{311.1, "E♭", 4},
		{659.3, "E", 5},
		{1397, "F", 6},
		{2960, "F#", 7},
		{6272, "G", 8},
		{25.96, "G#", 0},
		{55, "A", 1},
		{1865, "B♭", 6},
		{493.9, "B", 4},
	}

	for _, table := range tables {
		freqInfo := CreateFrequencyInfo(table.frequency)
		expected := FrequencyInfo{table.frequency, table.note, table.octave}
		if freqInfo != expected {
			t.Errorf("Expected: %v, Actual: %v", expected, freqInfo)
		}
	}
}

func TestZeroCrossing(t *testing.T) {
	tolerance := 5.0

	tables := []struct {
		size         int
		frequency    float64
		samplingRate float64
	}{
		{2048, 82.0, 48000.0},
		{2048, 82.0, 44100.0},
		{1024, 82.0, 48000.0},
		{2048, 110.0, 48000.0},
		{512, 440.0, 44100.0},
		{256, 1000.0, 44100.0},
	}

	for _, table := range tables {
		calculator := ZeroCrossingFrequencyCalculator{0.05}

		data := generateFrequency(table.size, table.frequency, table.samplingRate)
		actualFreq := calculator.GetFrequency(data, table.samplingRate)
		if math.Abs(actualFreq-table.frequency) >= tolerance {
			t.Errorf("Incorrect frequency when size = %d and sampling rate = %.2f, got: %.2f, expected: %.2f", table.size, table.samplingRate, actualFreq, table.frequency)
		}
	}
}

func generateFrequency(size int, frequency, samplingRate float64) []float64 {
	data := make([]float64, size)
	for i := range data {
		timectr := float64(i) / samplingRate
		data[i] = math.Sin(2 * math.Pi * frequency * timectr)
	}
	return data
}
