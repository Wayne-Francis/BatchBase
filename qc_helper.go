package main

type QCResult struct {
	Replicate int
	Result    float64
}

func addMultipleQCResults(results []float64) []QCResult {
	qcResults := []QCResult{}

	for i, result := range results {
		replicate := i + 1

		qcResults = append(qcResults, QCResult{
			Replicate: replicate,
			Result:    result,
		})
	}

	return qcResults
}
