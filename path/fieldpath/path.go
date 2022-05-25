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

package fieldpath

import (
	"encoding/json"
	"strings"
)

// Path provides a JSONPath-like struct and field-member "route" to a
// particular field within a resource. Path implements json.Marshaler
// interface.
type Path struct {
	parts []string
}

// String returns the dotted-notation representation of the Path
func (p *Path) String() string {
	return strings.Join(p.parts, ".")
}

// MarshalJSON returns the JSON encoding of a Path object.
func (p *Path) MarshalJSON() ([]byte, error) {
	// Since json.Marshal doesn't encode unexported struct fields we have to
	// copy the Path instance into a new struct object with exported fields.
	// See https://github.com/aws-controllers-k8s/community/issues/772
	return json.Marshal(
		struct {
			Parts []string
		}{
			p.parts,
		},
	)
}

// Pop removes the last part from the Path and returns it.
func (p *Path) Pop() (part string) {
	if len(p.parts) > 0 {
		part = p.parts[len(p.parts)-1]
		p.parts = p.parts[:len(p.parts)-1]
	}
	return part
}

// At returns the part of the Path at the supplied index, or empty string if
// index exceeds boundary.
func (p *Path) At(index int) string {
	if index < 0 || len(p.parts) == 0 || index > len(p.parts)-1 {
		return ""
	}
	return p.parts[index]
}

// Front returns the first part of the Path or empty string if the Path has no
// parts.
func (p *Path) Front() string {
	if len(p.parts) == 0 {
		return ""
	}
	return p.parts[0]
}

// PopFront removes the first part of the Path and returns it.
func (p *Path) PopFront() (part string) {
	if len(p.parts) > 0 {
		part = p.parts[0]
		p.parts = p.parts[1:]
	}
	return part
}

// Back returns the last part of the Path or empty string if the Path has no
// parts.
func (p *Path) Back() string {
	if len(p.parts) == 0 {
		return ""
	}
	return p.parts[len(p.parts)-1]
}

// PushBack adds a new part to the end of the Path.
func (p *Path) PushBack(part string) {
	p.parts = append(p.parts, part)
}

// Copy returns a new Path that is a copy of this Path
func (p *Path) Copy() *Path {
	return &Path{p.parts}
}

// CopyAt returns a new Path that is a copy of this Path up to the supplied
// index.
//
// e.g. given Path $A containing "X.Y", $A.CopyAt(0) would return a new Path
// containing just "X". $A.CopyAt(1) would return a new Path containing "X.Y".
func (p *Path) CopyAt(index int) *Path {
	if index < 0 || len(p.parts) == 0 || index > len(p.parts)-1 {
		return nil
	}
	return &Path{p.parts[0 : index+1]}
}

// Empty returns true if there are no parts to the Path
func (p *Path) Empty() bool {
	return len(p.parts) == 0
}

// Size returns the Path number of parts
func (p *Path) Size() int {
	return len(p.parts)
}

// HasPrefix returns true if the supplied string, delimited on ".", matches
// p.parts up to the length of the supplied string.
// e.g. if the Path p represents "A.B":
//  subject "A" -> true
//  subject "A.B" -> true
//  subject "A.B.C" -> false
//  subject "B" -> false
//  subject "A.C" -> false
func (p *Path) HasPrefix(subject string) bool {
	subjectSplit := strings.Split(subject, ".")

	if len(subjectSplit) > len(p.parts) {
		return false
	}

	for i, s := range subjectSplit {
		if p.parts[i] != s {
			return false
		}
	}

	return true
}

// HasPrefixFold is the same as HasPrefix but uses case-insensitive comparisons
func (p *Path) HasPrefixFold(subject string) bool {
	subjectSplit := strings.Split(subject, ".")

	if len(subjectSplit) > len(p.parts) {
		return false
	}

	for i, s := range subjectSplit {
		if !strings.EqualFold(p.parts[i], s) {
			return false
		}
	}

	return true
}

// FromString returns a new Path from a dotted-notation string, e.g.
// "Author.Name".
func FromString(dotted string) *Path {
	return &Path{strings.Split(dotted, ".")}
}
