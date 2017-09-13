//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package markov

import (
	"errors"
	"fmt"
	"math/rand"
)

//-----------------------------------------------------------------------------

// Return the interval corresponding to the x value.
func locate(x float64, p []float64) int {
	if x < 0 || x > 1 {
		panic("x is out of range")
	}
	// TODO use bisecting search
	for i := 0; true; i += 1 {
		x -= p[i]
		if x < 0 {
			return i
		}
	}
	panic("interval not found")
}

// Normalize the probability values of a slice of floats.
func normalize(x []float64) error {
	if x == nil {
		return errors.New("nil value")
	}
	n := len(x)
	if n == 0 {
		return errors.New("no values")
	}
	sum := 0.0
	for i := range x {
		if x[i] < 0 {
			return errors.New("value < 0")
		}
		sum += x[i]
	}
	for i := range x {
		x[i] = x[i] / sum
	}
	return nil
}

//-----------------------------------------------------------------------------

// Markov Model
type MM struct {
	state int       // current state
	n     int       // number of states
	a     []float64 // state transition probabilities
	pi    []float64 // initial state probabilties
}

func NewMM(
	n int, // number of states
	a []float64, // state transition probabilities
	pi []float64, // initial state probabilties
) (*MM, error) {

	// number of states
	if n <= 0 {
		return nil, errors.New("bad number of states")
	}

	// state transition matrix
	if a == nil {
		a = make([]float64, n*n)
		for i := range a {
			a[i] = 1
		}
	}
	if len(a) != n*n {
		return nil, fmt.Errorf("state transition matrix must have %d elements", n*n)
	}
	for i := 0; i < n; i++ {
		err := normalize(a[i*n : (i+1)*n])
		if err != nil {
			return nil, err
		}
	}

	// initial state probabilities
	if pi == nil {
		pi = make([]float64, n)
		for i := range pi {
			pi[i] = 1
		}
	}
	if len(pi) != n {
		return nil, fmt.Errorf("initial state matrix must have %d elements", n)
	}
	normalize(pi)

	return &MM{n: n, a: a, pi: pi}, nil
}

// Initialise the markov model, return the initial state.
func (mm *MM) Init() int {
	mm.state = locate(rand.Float64(), mm.pi)
	return mm.state
}

// Return the current state.
func (mm *MM) State() int {
	return mm.state
}

// Transition to the next state.
func (mm *MM) Next() int {
	i := mm.state * mm.n
	mm.state = locate(rand.Float64(), mm.a[i:i+mm.n])
	return mm.state
}

//-----------------------------------------------------------------------------

// Hidden Markov Model
type HMM struct {
	mm *MM       // markov model
	s  int       // number of output symbols
	b  []float64 // symbol emission probabilities
}

func NewHMM(
	n int, // number of states
	a []float64, // state transition probabilities
	pi []float64, // initial state probabilties
	s int, // number of output symbols
	b []float64, // symbol emission probabilities
) (*HMM, error) {

	// markov model
	mm, err := NewMM(n, a, pi)
	if err != nil {
		return nil, err
	}

	// number of symbols
	if s <= 0 {
		return nil, errors.New("bad number of symbols")
	}

	// symbol emission matrix
	if b == nil {
		b = make([]float64, n*s)
		for i := range b {
			b[i] = 1
		}
	}
	if len(b) != n*s {
		return nil, fmt.Errorf("symbol emission matrix must have %d elements", n*s)
	}
	for i := 0; i < n; i++ {
		err := normalize(b[i*s : (i+1)*s])
		if err != nil {
			return nil, err
		}
	}

	return &HMM{mm: mm, s: s, b: b}, nil
}

//-----------------------------------------------------------------------------
