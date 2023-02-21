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

// Scalar represent a set of useful types to be used with type parameters.
//
// For now we only support the type that are frequently used within the
// controller repositories, code-generator and runtime.
type Scalar interface {
	int | int32 | int64 | float32 | float64 | bool | string
}

// EqualPStrict returns true if two given scalar pointers are strictly
// equal, meaning that they either have the same nility or the same value.
func EqualPStrict[s Scalar](x, y *s) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return false
	}
	return *x == *y
}

// EqualP returns true if two given scalars represent the same value. It
// behaves similarilyu to EqualPStrict but will also return true if one of
// the parameters is nil and the other points to zero scalar value.
func EqualP[s Scalar](x, y *s) bool {
	var zero s
	if x == nil {
		return y == nil || *y == zero
	}
	if y == nil {
		return *x == zero
	}
	return *x == *y
}
