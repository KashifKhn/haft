package detector

const (
	ConfidenceHigh      = 0.85
	ConfidenceThreshold = 0.70
	ConfidenceMedium    = 0.50
	ConfidenceLow       = 0.30
)

const (
	WeightSignalStrength = 0.40
	WeightSampleSize     = 0.30
	WeightConsistency    = 0.30
)

const (
	MinSampleSize       = 3
	OptimalSampleSize   = 10
	MaxSampleMultiplier = 1.0
)

type ConfidenceCalculator struct{}

func NewConfidenceCalculator() *ConfidenceCalculator {
	return &ConfidenceCalculator{}
}

func (c *ConfidenceCalculator) Calculate(signalStrength float64, sampleSize int, consistency float64) float64 {
	normalizedSampleSize := c.normalizeSampleSize(sampleSize)

	confidence := (signalStrength * WeightSignalStrength) +
		(normalizedSampleSize * WeightSampleSize) +
		(consistency * WeightConsistency)

	return clamp(confidence, 0.0, 1.0)
}

func (c *ConfidenceCalculator) CalculateFromCounts(matchCount, totalCount int) float64 {
	if totalCount == 0 {
		return 0.0
	}

	signalStrength := float64(matchCount) / float64(totalCount)
	consistency := c.calculateConsistencyFromRatio(signalStrength)

	return c.Calculate(signalStrength, totalCount, consistency)
}

func (c *ConfidenceCalculator) CalculateArchitectureConfidence(
	matchingFiles, totalFiles int,
	hasDistinctiveMarkers bool,
	ambiguityScore float64,
) float64 {
	if totalFiles == 0 {
		return 0.0
	}

	baseRatio := float64(matchingFiles) / float64(totalFiles)

	markerBonus := 0.0
	if hasDistinctiveMarkers {
		markerBonus = 0.10
	}

	ambiguityPenalty := ambiguityScore * 0.15

	confidence := c.Calculate(baseRatio, totalFiles, 1.0-ambiguityScore)
	confidence = confidence + markerBonus - ambiguityPenalty

	return clamp(confidence, 0.0, 1.0)
}

func (c *ConfidenceCalculator) CalculatePatternConfidence(
	occurrences int,
	sampleSize int,
	patternVariations int,
) float64 {
	if sampleSize == 0 {
		return 0.0
	}

	signalStrength := float64(occurrences) / float64(sampleSize)

	consistency := 1.0
	if patternVariations > 1 {
		consistency = 1.0 / float64(patternVariations)
	}

	return c.Calculate(signalStrength, sampleSize, consistency)
}

func (c *ConfidenceCalculator) normalizeSampleSize(sampleSize int) float64 {
	if sampleSize <= 0 {
		return 0.0
	}

	if sampleSize >= OptimalSampleSize {
		return MaxSampleMultiplier
	}

	return float64(sampleSize) / float64(OptimalSampleSize)
}

func (c *ConfidenceCalculator) calculateConsistencyFromRatio(ratio float64) float64 {
	if ratio >= 0.9 {
		return 1.0
	}
	if ratio >= 0.7 {
		return 0.8
	}
	if ratio >= 0.5 {
		return 0.6
	}
	return 0.4
}

func (c *ConfidenceCalculator) IsHighConfidence(confidence float64) bool {
	return confidence >= ConfidenceHigh
}

func (c *ConfidenceCalculator) MeetsThreshold(confidence float64) bool {
	return confidence >= ConfidenceThreshold
}

func (c *ConfidenceCalculator) NeedsUserConfirmation(confidence float64) bool {
	return confidence < ConfidenceThreshold && confidence >= ConfidenceLow
}

func (c *ConfidenceCalculator) IsTooLow(confidence float64) bool {
	return confidence < ConfidenceLow
}

func (c *ConfidenceCalculator) CompareDetections(first, second float64) int {
	diff := first - second

	if diff > 0.1 {
		return 1
	}
	if diff < -0.1 {
		return -1
	}
	return 0
}

func (c *ConfidenceCalculator) CombineConfidences(confidences ...float64) float64 {
	if len(confidences) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, conf := range confidences {
		sum += conf
	}

	return sum / float64(len(confidences))
}

func (c *ConfidenceCalculator) WeightedCombine(values []float64, weights []float64) float64 {
	if len(values) == 0 || len(values) != len(weights) {
		return 0.0
	}

	totalWeight := 0.0
	weightedSum := 0.0

	for i, value := range values {
		weightedSum += value * weights[i]
		totalWeight += weights[i]
	}

	if totalWeight == 0 {
		return 0.0
	}

	return weightedSum / totalWeight
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func (c *ConfidenceCalculator) GetConfidenceLevel(confidence float64) string {
	switch {
	case confidence >= ConfidenceHigh:
		return "high"
	case confidence >= ConfidenceThreshold:
		return "medium"
	case confidence >= ConfidenceLow:
		return "low"
	default:
		return "very_low"
	}
}

func (c *ConfidenceCalculator) FormatPercentage(confidence float64) int {
	return int(confidence * 100)
}
