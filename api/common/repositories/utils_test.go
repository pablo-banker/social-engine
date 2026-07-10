package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_isEmptyID(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		assert.True(t, isEmptyID(nil))
	})

	t.Run("zero uuid value", func(t *testing.T) {
		var id uuid.UUID
		assert.True(t, isEmptyID(id))
	})

	t.Run("non-zero uuid value", func(t *testing.T) {
		id := uuid.New()
		assert.False(t, isEmptyID(id))
	})

	t.Run("nil uuid pointer", func(t *testing.T) {
		var id *uuid.UUID
		assert.True(t, isEmptyID(id))
	})

	t.Run("pointer to zero uuid", func(t *testing.T) {
		var v uuid.UUID
		id := &v
		assert.True(t, isEmptyID(id))
	})

	t.Run("pointer to non-zero uuid", func(t *testing.T) {
		v := uuid.New()
		id := &v
		assert.False(t, isEmptyID(id))
	})

	t.Run("zero int", func(t *testing.T) {
		var n int
		assert.True(t, isEmptyID(n))
	})

	t.Run("non-zero int", func(t *testing.T) {
		n := 42
		assert.False(t, isEmptyID(n))
	})

	t.Run("empty string", func(t *testing.T) {
		s := ""
		assert.True(t, isEmptyID(s))
	})

	t.Run("non-empty string", func(t *testing.T) {
		s := "abc"
		assert.False(t, isEmptyID(s))
	})

	t.Run("zero struct value", func(t *testing.T) {
		type dummy struct {
			A int
		}
		var d dummy
		assert.True(t, isEmptyID(d))
	})

	t.Run("non-zero struct value", func(t *testing.T) {
		type dummy struct {
			A int
		}
		d := dummy{A: 1}
		assert.False(t, isEmptyID(d))
	})

	t.Run("pointer to zero struct", func(t *testing.T) {
		type dummy struct {
			A int
		}
		var d dummy
		assert.True(t, isEmptyID(&d))
	})

	t.Run("pointer to non-zero struct", func(t *testing.T) {
		type dummy struct {
			A int
		}
		d := dummy{A: 10}
		assert.False(t, isEmptyID(&d))
	})
}

func Test_Normalize(t *testing.T) {
	t.Run("trim spaces only", func(t *testing.T) {
		in := "   hello world   "
		out := Normalize(in)
		assert.Equal(t, "hello world", out)
	})

	t.Run("remove tabs only", func(t *testing.T) {
		in := "\thello\tworld\t"
		out := Normalize(in)
		// tabs viram nada, mas sem espaços extras
		assert.Equal(t, "helloworld", out)
	})

	t.Run("spaces and tabs mixed", func(t *testing.T) {
		in := " \t  hello \t world \t "
		out := Normalize(in)
		// tabs são removidos, depois trim
		assert.Equal(t, "hello  world", out)
	})

	t.Run("no change", func(t *testing.T) {
		in := "already-normal"
		out := Normalize(in)
		assert.Equal(t, "already-normal", out)
	})

	t.Run("only whitespace", func(t *testing.T) {
		in := " \t \t  "
		out := Normalize(in)
		assert.Equal(t, "", out)
	})
}
