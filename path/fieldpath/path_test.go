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

package fieldpath_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/pkg/path/fieldpath"
)

func TestBasics(t *testing.T) {
	require := require.New(t)

	pname := fieldpath.FromString("Author.Name")
	require.Equal("Author.Name", pname.String())

	pstate := fieldpath.FromString("Author.Address.State")
	require.Equal("Author.Address.State", pstate.String())

	require.Equal("Author", pstate.Front())
	require.Equal("State", pstate.Back())

	require.Equal("Author", pstate.At(0))
	require.Equal("Address", pstate.At(1))
	require.Equal("State", pstate.At(2))
	require.Equal("", pstate.At(3))

	pauth := pstate.CopyAt(0)
	require.Equal("Author", pauth.String())

	last := pstate.Pop()
	require.Equal("State", last)
	require.Equal("Address", pstate.Back())

	pstate.PushBack("Country")
	require.Equal("Country", pstate.Back())

	front := pstate.PopFront()
	require.Equal("Author", front)
	require.Equal("Address", pstate.Front())
	require.False(pstate.Empty())
	pstate.Pop()
	require.False(pstate.Empty())
	pstate.Pop()
	require.True(pstate.Empty())
}

func TestHasPrefix(t *testing.T) {
	require := require.New(t)

	p := fieldpath.FromString("Author.Name")
	require.True(p.HasPrefix("Author.Name"))
	require.True(p.HasPrefix("Author"))
	require.False(p.HasPrefix("Name"))
	require.False(p.HasPrefix("Author.Address"))
	// Case-insensitive comparisons...
	require.False(p.HasPrefix("author"))
	require.True(p.HasPrefixFold("author"))
}
