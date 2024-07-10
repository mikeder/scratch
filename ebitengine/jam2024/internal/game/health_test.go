package game_test

import (
	"testing"

	"github.com/mikeder/scratchygo/ebitengine/jam2024/internal/game"
)

func TestHealth(t *testing.T) {
	t.Run("add", func(t *testing.T) {

		have := game.Health(20)
		have.Add(10)
		want := game.Health(30)

		if have != want {
			t.Logf("have: %d, want: %d", have, want)
			t.FailNow()
		}
	})
	t.Run("remove", func(t *testing.T) {

		have := game.Health(20)
		have.Remove(10)
		want := game.Health(10)

		if have != want {
			t.Logf("have: %d, want: %d", have, want)
			t.FailNow()
		}
	})
	t.Run("remove to 0", func(t *testing.T) {

		have := game.Health(20)
		have.Remove(30)
		want := game.Health(0)

		if have != want {
			t.Logf("have: %d, want: %d", have, want)
			t.FailNow()
		}
	})

}
