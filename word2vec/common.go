package word2vec

type WordDistance struct {
	WordA string
	WordB string
	Distance float64
}

type Distances []WordDistance

func (slice Distances) Len() int {
	return len(slice)
}

func (slice Distances) Less(i, j int) bool {
	return slice[i].Distance < slice[j].Distance;
}

func (slice Distances) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
