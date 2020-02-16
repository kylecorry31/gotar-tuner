package gotar

// FrequencyCalculator a frequency calculator
type FrequencyCalculator interface {
	GetFrequency(waveform []float64, sampleRate float64) float64
}
