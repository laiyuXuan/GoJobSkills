package word2vec

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"time"
)

func LoadModel(path string) map[string][]float64 {
	before := time.Now().Second()
	file, _ := os.Open(path)

	defer file.Close()
	vectorMap := make(map[string][]float64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		word := strings.Split(line, " ")[0]
		split := strings.Split(strings.TrimSpace(line[1:]), " ")
		vectors := make([]float64, 0)
		for _, vector := range split {
			float, _ := strconv.ParseFloat(vector, 64)
			vectors = append(vectors, float)
		}
		vectorMap[word] = vectors
	}
	after := time.Now().Second()
	logger.Printf("load word2vec model took %d seconds", after - before)
	return vectorMap
}
