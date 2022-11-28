/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2022, InterlockLedger Network
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * * Neither the name of the copyright holder nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package testbinary

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextSeed(t *testing.T) {
	queue := make(chan int64)

	// Start at lest 4 concurrent consumers
	for i := 0; i < 4; i++ {
		go func(queue chan int64) {
			for j := 0; j < 10; j++ {
				queue <- NextSeed()
			}
		}(queue)
	}

	// Ensure the the concurrent calls will never share the same value
	values := make(map[int64]bool)
	for i := 0; i < 40; i++ {
		v := <-queue
		_, found := values[v]
		assert.False(t, found)
		values[v] = true
	}
}

func TestNewSeededRandom(t *testing.T) {

	r1 := NewSeededRandom()
	r2 := NewSeededRandom()

	dup := 0
	for i := 0; i < 100; i++ {
		if r1.Int() == r2.Int() {
			dup++
		}
	}
	// The sequences must have no more that 10% of random values in the same
	// order. This criteria is completely arbitrary because the idea is to prove
	// that both r1 and r2 have distinct seed.
	assert.Greater(t, 10, dup)
}

func TestFillRandomWithSeed(t *testing.T) {

	for size := 1; size < 50; size++ {
		b1 := make([]byte, size)
		b2 := make([]byte, size)
		b3 := make([]byte, size)

		c1 := FillRandomWithSeed(b1, 1234)
		assert.Same(t, &b1[0], &c1[0])
		c2 := FillRandomWithSeed(b2, 1234)
		assert.Same(t, &b2[0], &c2[0])
		assert.Equal(t, c1, c2)

		c3 := FillRandomWithSeed(b3, 1235)
		assert.Same(t, &b3[0], &c3[0])
		assert.NotEqual(t, c1, c3)
	}

	// nil
	assert.Nil(t, FillRandomWithSeed(nil, 10))
	b := make([]byte, 1)
	// empty
	c := FillRandomWithSeed(b[:0], 10)
	assert.Empty(t, c)
	c2 := c[0:cap(c):cap(c)]
	assert.Same(t, &b[0], &c2[0])
}

func TestFillRandom(t *testing.T) {

	for size := 1; size < 50; size++ {
		b1 := make([]byte, size)
		b2 := make([]byte, size)

		c1 := FillRandom(b1)
		assert.Same(t, &b1[0], &c1[0])
		c2 := FillRandom(b2)
		assert.Same(t, &b2[0], &c2[0])
		assert.NotEqual(t, c1, c2)
	}

	// nil
	assert.Nil(t, FillRandom(nil))
	b := make([]byte, 1)
	// empty
	c := FillRandom(b[:0])
	assert.Empty(t, c)
	c2 := c[0:cap(c):cap(c)]
	assert.Same(t, &b[0], &c2[0])
}

func TestFillSeq(t *testing.T) {

	for size := 1; size < 50; size++ {
		b1 := make([]byte, size)
		c1 := FillSeq(b1)
		assert.Same(t, &b1[0], &c1[0])
		for j, v := range b1 {
			assert.Equal(t, byte(j), v)
		}
	}

	// nil
	assert.Nil(t, FillSeq(nil))
	b := make([]byte, 1)
	// empty
	c := FillSeq(b[:0])
	c2 := c[0:cap(c):cap(c)]
	assert.Same(t, &b[0], &c2[0])
}

func TestFillSeqCustom(t *testing.T) {

	for size := 1; size < 50; size++ {
		b1 := make([]byte, size)
		c1 := FillSeqCustom(b1, size)
		assert.Same(t, &b1[0], &c1[0])
		for j, v := range b1 {
			assert.Equal(t, byte(j+size), v)
		}
	}

	// nil
	assert.Nil(t, FillSeqCustom(nil, 1))
	b := make([]byte, 1)
	// empty
	c := FillSeqCustom(b[:0], 1)
	c2 := c[0:cap(c):cap(c)]
	assert.Same(t, &b[0], &c2[0])
}

func TestNewRandomString(t *testing.T) {

	for size := 1; size < 50; size++ {
		// Generate a new string.
		s1 := NewRandomString(size)
		assert.Equal(t, size, utf8.RuneCountInString(s1))
		// Check if all runes are graphic
		r := strings.NewReader(s1)
		for r.Len() > 0 {
			c, _, err := r.ReadRune()
			require.Nil(t, err)
			assert.True(t, unicode.IsGraphic(c))
		}
		s2 := NewRandomString(size)
		assert.Equal(t, size, utf8.RuneCountInString(s2))
		assert.NotEqual(t, s1, s2)
	}
}
