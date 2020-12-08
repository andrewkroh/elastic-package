// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package multierror

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnique(t *testing.T) {
	errs := Error{
		errors.New("2"),
		errors.New("1"),
		errors.New("2"),
		errors.New("1"),
		errors.New("3"),
	}

	unique := errs.Unique()

	require.Len(t, unique, 3)
	require.Len(t, errs, 5)
}

func TestError(t *testing.T) {
	errs := Error{
		errors.New("one"),
		errors.New("two"),
		errors.New("three"),
		errors.New("four"),
		errors.New("five"),
	}

	const expected = `[0] one
[1] two
[2] three
[3] four
[4] five`

	assert.Equal(t, expected, errs.Error())
}

func BenchmarkError_Unique(b *testing.B) {
	var err Error
	for i := 0; i < 500; i++ {
		err = append(err, fmt.Errorf("%d: %s", i, strings.Repeat("hello", 40)))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err.Unique()
	}
}

func BenchmarkError_String(b *testing.B) {
	var err Error
	for i := 0; i < 500; i++ {
		err = append(err, fmt.Errorf("%d: %s", i, strings.Repeat("hello", 40)))
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err.Error()
	}
}
