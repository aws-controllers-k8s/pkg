// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package scalar

import (
	"testing"
)

// Temporary solution before bringing samber/lo add a dependency
// For now this is only needed for uni tests.
func ptr[T any](x T) *T {
	return &x
}

func TestEqualScalarPStrict_int(t *testing.T) {
	type args[scalar Scalar] struct {
		x *scalar
		y *scalar
	}

	type intArgs args[int]
	tests := []struct {
		name string
		args intArgs
		want bool
	}{
		{
			name: "nil pointers",
			args: intArgs{nil, nil},
			want: true,
		},
		{
			name: "one nil pointer - left, with zero value - right",
			args: intArgs{nil, ptr(0)},
			want: false,
		},
		{
			name: "one nil pointer - right, with zero value - left",
			args: intArgs{ptr(0), nil},
			want: false,
		},
		{
			name: "non nil pointers - equal values",
			args: intArgs{ptr(128), ptr(128)},
			want: true,
		},
		{
			name: "non nil pointers - non equal values",
			args: intArgs{ptr(128), ptr(64)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualPStrict(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("EqualScalarPStrict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualScalarP_int(t *testing.T) {
	type args[scalar Scalar] struct {
		x *scalar
		y *scalar
	}
	type intArgs args[int]
	tests := []struct {
		name string
		args intArgs
		want bool
	}{
		{
			name: "nil pointers",
			args: intArgs{nil, nil},
			want: true,
		},
		{
			name: "one nil pointer - left, with zero value - right",
			args: intArgs{nil, ptr(0)},
			want: true,
		},
		{
			name: "one nil pointer - right, with zero value - left",
			args: intArgs{ptr(0), nil},
			want: true,
		},
		{
			name: "zero values",
			args: intArgs{ptr(0), ptr(0)},
			want: true,
		},
		{
			name: "non nil pointers - equal",
			args: intArgs{ptr(128), ptr(128)},
			want: true,
		},
		{
			name: "non nil pointers - not equal",
			args: intArgs{ptr(128), ptr(64)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualP(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("EqualScalarP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualScalarPStrict_string(t *testing.T) {
	type args[scalar Scalar] struct {
		x *scalar
		y *scalar
	}
	type stringArgs args[string]
	tests := []struct {
		name string
		args stringArgs
		want bool
	}{
		{
			name: "nil pointers",
			args: stringArgs{nil, nil},
			want: true,
		},
		{
			name: "one nil pointer - left, with zero value - right",
			args: stringArgs{nil, ptr("")},
			want: false,
		},
		{
			name: "one nil pointer - right, with zero value - left",
			args: stringArgs{ptr(""), nil},
			want: false,
		},
		{
			name: "non nil pointers - equal values",
			args: stringArgs{ptr("ABC"), ptr("ABC")},
			want: true,
		},
		{
			name: "non nil pointers - non equal values",
			args: stringArgs{ptr("ABC"), ptr("XYZ")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualPStrict(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("EqualScalarPStrict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualScalarP_string(t *testing.T) {
	type args[scalar Scalar] struct {
		x *scalar
		y *scalar
	}
	type stringArgs args[string]
	tests := []struct {
		name string
		args stringArgs
		want bool
	}{
		{
			name: "nil pointers",
			args: stringArgs{nil, nil},
			want: true,
		},
		{
			name: "one nil pointer - left, with zero value - right",
			args: stringArgs{nil, ptr("")},
			want: true,
		},
		{
			name: "one nil pointer - right, with zero value - left",
			args: stringArgs{ptr(""), nil},
			want: true,
		},
		{
			name: "zero values",
			args: stringArgs{ptr(""), ptr("")},
			want: true,
		},
		{
			name: "non nil pointers - equal",
			args: stringArgs{ptr("ABC"), ptr("ABC")},
			want: true,
		},
		{
			name: "non nil pointers - not equal",
			args: stringArgs{ptr("ABC"), ptr("XYZ")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualP(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("EqualScalarP() = %v, want %v", got, tt.want)
			}
		})
	}
}
