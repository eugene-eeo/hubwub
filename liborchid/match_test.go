package liborchid_test

import "fmt"
import "testing"
import "github.com/eugene-eeo/orchid/liborchid"
import "github.com/stretchr/testify/assert"

type matchTest struct {
	query    string
	haystack string
	match    bool
	distance int
}

func TestMatch(t *testing.T) {
	tests := []matchTest{
		matchTest{"def", "define", true, 0},
		matchTest{"def", "deefine", true, 1},
		matchTest{"d", "efine", false, -1},
		matchTest{"mid", "start mid end", true, 0},
		matchTest{"谢谢", "CJK谢谢", true, 0},
	}

	for _, m := range tests {
		matched, distance := liborchid.Match(m.query, m.haystack)
		assert.Equal(
			t,
			m.match,
			matched,
			fmt.Sprintf("%s : %s", m.query, m.haystack),
		)
		if m.match {
			assert.Equal(
				t,
				m.distance,
				distance,
				fmt.Sprintf("%s : %s", m.query, m.haystack),
			)
		}
	}
}
