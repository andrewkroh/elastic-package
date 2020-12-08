// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package multierror

import (
	"sort"
	"strconv"
	"strings"
)

// Error is a multi-error representation.
type Error []error

// Unique selects only unique errors based on their string value.
func (me Error) Unique() Error {
	return deduplicate(me)
}

// deduplicate a list of errors based on their string value. The returned
// list of errors is sorted.
func deduplicate(errs []error) []error {
	// errorString contains a cached copy of the error's message so that
	// each error is only converted to a string one time.
	type errorString struct {
		err error
		msg string
	}

	// deduplicate
	m := make(map[string]errorString, len(errs))
	for _, e := range errs {
		es := errorString{err: e, msg: e.Error()}
		m[es.msg] = es
	}

	// sort
	tmp := make([]errorString, 0, len(m))
	for _, es := range m {
		tmp = append(tmp, es)
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].msg < tmp[j].msg
	})

	// make errors slice
	out := make([]error, 0, len(tmp))
	for _, es := range tmp {
		out = append(out, es.err)
	}
	return out
}

// Error combines a detailed report consisting of attached errors separated with new lines.
func (me Error) Error() string {
	if me == nil {
		return ""
	}

	sb := new(strings.Builder)
	for i, err := range me {
		sb.WriteString("[")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString("] ")
		sb.WriteString(err.Error())
		if i < len(me)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
