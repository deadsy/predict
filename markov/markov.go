//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package markov

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
)

//-----------------------------------------------------------------------------

const EPSILON = 1e-12

//-----------------------------------------------------------------------------

func FloatDecode(x float64) string {
	i := math.Float64bits(x)
	s := int((i >> 63) & 1)
	f := i & ((1 << 52) - 1)
	e := int((i>>52)&((1<<11)-1)) - 1023
	return fmt.Sprintf("s %d f 0x%013x e %d", s, f, e)
}

const min_normal = 2.2250738585072014E-308 // 2**-1022

func EqualFloat64(a, b, epsilon float64) bool {

	log.Printf("a = %s b = %s", FloatDecode(a), FloatDecode(b))

	if a == b {
		return true
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)
	if a == 0 || b == 0 || diff < min_normal {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * min_normal)
	}
	// use relative error
	return diff/math.Min((absA+absB), math.MaxFloat64) < epsilon
}

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

type MarkovModel struct {
	state int       // current state
	n     int       // number of states
	a     []float32 // state transition probabilities
	pi    []float32 // initial state probabilties
}

func NewMarkovModel(
	n int, // number of states
	a []float32, // state transition probabilities
	pi []float32, // initial state probabilties
) (*MarkovModel, error) {
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
		x := 0.0
		for j := 0; j < n; j++ {
			x += float64(a[(i*n)+j])
		}
		if !EqualFloat64(x, 1.0, EPSILON) {
			return nil, fmt.Errorf("state transition matrix row %d does not sum to 1.0", i)
		}
	}
	// check the initial state probabilities
	x := 0.0
	for i := 0; i < n; i++ {
		x += float64(pi[i])
	}
	if x != 1.0 {
		return nil, errors.New("initial state matrix does not sum to 1.0")
	}
	return &MarkovModel{n: n, a: a, pi: pi}, nil
}

//-----------------------------------------------------------------------------

// Initialise the markov model, return the initial state.
func (mm *MarkovModel) Init() int {
	mm.state = locate(rand.Float32(), mm.pi)
	return mm.state
}

// Return the current state.
func (mm *MarkovModel) State() int {
	return mm.state
}

// Transition to the next state.
func (mm *MarkovModel) Next() int {
	i := mm.state * mm.n
	mm.state = locate(rand.Float32(), mm.a[i:i+mm.n])
	return mm.state
}

//-----------------------------------------------------------------------------
