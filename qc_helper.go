package main

type QCResult struct {
	Replicate int
	Result    float64
}

func addMultipleQCResults(results []float64, startingReplicate int) []QCResult {
	qcResults := []QCResult{}

	for i, result := range results {
		qcResults = append(qcResults, QCResult{
			Replicate: startingReplicate + i,
			Result:    result,
		})
	}

	return qcResults
}
