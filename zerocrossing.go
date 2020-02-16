package gotar

// ZeroCrossingFrequencyCalculator a frequency calculator which counts the number of zero crossings
type ZeroCrossingFrequencyCalculator struct {
	CrossingThreshold float64
}

// GetFrequency get the primary frequency of a waveform
func (c ZeroCrossingFrequencyCalculator) GetFrequency(waveform []float64, sampleRate float64) float64 {
	positiveThresh := c.CrossingThreshold
	negativeThresh := -positiveThresh

	isPositive := waveform[0] > 0

	firstChange := -1
	secondChange := -1

	changeCnt := 0

	for i, value := range waveform {
		wasPositive := isPositive
		isPositive = shouldBeOn(isPositive, value, positiveThresh, negativeThresh)

		if wasPositive && !isPositive {
			changeCnt++

			if changeCnt == 1 {
				firstChange = i
			} else if changeCnt == 2 {
				secondChange = i
				break
			}
		}
	}

	if firstChange != -1 && secondChange != -1 {
		totalTime := float64(len(waveform)) / sampleRate
		period := totalTime * float64(secondChange-firstChange) / (float64(len(waveform)))
		return 1 / period
	}

	return 0
}

func shouldBeOn(isOn bool, currentValue float64, onThreshold float64, offThreshold float64) bool {
	if isOn && currentValue <= offThreshold {
		return false
	} else if !isOn && currentValue >= onThreshold {
		return true
	}

	return isOn
}
