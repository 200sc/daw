package alg

type stwHeap struct {
	bh           []float64
	weightsBelow []float64
}

// Select Total Weight Heap
// This name was chosen relatively arbitrarily, if there
// is a canonical academic name for this structure we'd gladly
// use that instead
func newSTWHeap(f []float64) *stwHeap {
	stwh := new(stwHeap)
	f = append([]float64{0}, f...)
	// The order of elements literally does not
	// matter, so 'heap' is a misnomer.
	stwh.bh = f
	stwh.weightsBelow = make([]float64, len(f))
	copy(stwh.weightsBelow, f)
	for i := len(f) - 1; i > 1; i-- {
		stwh.weightsBelow[i>>1] += stwh.weightsBelow[i]
	}
	return stwh
}

func (stwh *stwHeap) Pop(rng float64) int {
	if stwh.weightsBelow[1] <= ε {
		return -1
	}
	w := stwh.weightsBelow[1] * rng
	i := 1

	// With the >= here, we don't accept 0 weights
	for w >= stwh.bh[i] {
		w -= stwh.bh[i]
		i <<= 1 // Propagate to left child
		if w >= stwh.weightsBelow[i] {
			w -= stwh.weightsBelow[i]
			i++ // Switch to right child
		}
	}
	i2 := i
	w = stwh.bh[i]
	// Instead of removing a node we set its weight to 0.
	stwh.bh[i] = 0

	// All parents of the index need to be reduced
	// in total weight.
	for i > 0 {
		stwh.weightsBelow[i] -= w
		i >>= 1
	}
	return i2 - 1
}
