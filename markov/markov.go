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
func locate(x float32, p []float32) int {
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

//-----------------------------------------------------------------------------

// Markov Model
type MM struct {
	state int       // current state
	n     int       // number of states
	a     []float32 // state transition probabilities
	pi    []float32 // initial state probabilties
}

func NewMM(
	n int, // number of states
	a []float32, // state transition probabilities
	pi []float32, // initial state probabilties
) (*MM, error) {
	// parameter checking
	if n <= 0 {
		return nil, errors.New("bad number of states")
	}
	if len(a) != n*n {
		return nil, fmt.Errorf("state transition matrix must have %d elements", n*n)
	}
	if len(pi) != n {
		return nil, fmt.Errorf("initial state matrix must have %d elements", n)
	}
	// check the state transition probabilities
	for i := 0; i < n; i++ {
		x := float32(0)
		for j := 0; j < n; j++ {
			x += a[(i*n)+j]
		}
		if x != 1 {
			return nil, fmt.Errorf("state transition matrix row %d does not sum to 1.0", i)
		}
	}
	// check the initial state probabilities
	x := float32(0)
	for i := 0; i < n; i++ {
		x += pi[i]
	}
	if x != 1 {
		return nil, errors.New("initial state matrix does not sum to 1.0")
	}
	return &MM{n: n, a: a, pi: pi}, nil
}

//-----------------------------------------------------------------------------

// Initialise the markov model, return the initial state.
func (mm *MM) Init() int {
	mm.state = locate(rand.Float32(), mm.pi)
	return mm.state
}

// Return the current state.
func (mm *MM) State() int {
	return mm.state
}

// Transition to the next state.
func (mm *MM) Next() int {
	i := mm.state * mm.n
	mm.state = locate(rand.Float32(), mm.a[i:i+mm.n])
	return mm.state
}

//-----------------------------------------------------------------------------
