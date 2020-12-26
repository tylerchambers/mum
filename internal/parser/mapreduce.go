package parser

// FindWordsWithMapReduce does a simple mapreduce to produce a for each file in all directories recursively.
func FindWordsWithMapReduce(directory string, poolSize int) map[string]int {

	reduceInput := make(chan map[string]int)
	reduceOutput := make(chan map[string]int)
	workerOutput := make(chan chan map[string]int, poolSize)

	go reducer(reduceInput, reduceOutput)

	go func() {
		for workerChan := range workerOutput {
			reduceInput <- <-workerChan
		}

		close(reduceInput)
	}()

	go func() {
		for item := range findFiles(directory) {
			myChan := make(chan map[string]int)
			go wordCount(item, myChan)
			workerOutput <- myChan
		}

		close(workerOutput)
	}()

	return <-reduceOutput
}
