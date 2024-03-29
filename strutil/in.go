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

package strutil

// InStrings returns true if the subject string is contained in the supplied
// slice of strings
func InStrings(subject string, collection []string) bool {
	for _, item := range collection {
		if subject == item {
			return true
		}
	}
	return false
}

// InStringPs returns true if the subject string is contained in the supplied
// slice of string pointers
func InStringPs(subject string, collection []*string) bool {
	for _, item := range collection {
		if subject == *item {
			return true
		}
	}
	return false
}
