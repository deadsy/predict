//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------

package markov

import (
	"log"
	"testing"
)

//-----------------------------------------------------------------------------

func Test_Markov(t *testing.T) {
	a := []float32{
		0.3, 0.7,
		0.4, 0.6,
	}
	pi := []float32{
		0.5, 0.5,
	}
	mm, err := NewMarkovModel(2, a, pi)
	if err != nil {
		t.Error(err)
		return
	}
	s := mm.Init()
	for i := 1; i < 10; i++ {
		log.Printf("%d", s)
		s = mm.Next()
	}
}

//-----------------------------------------------------------------------------
