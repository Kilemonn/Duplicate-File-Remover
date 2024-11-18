package argument_list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestArgumentList test different functions of the ArgumentList object, [flag.Value.String] and [flag.Value.Set].
func TestArgumentList(t *testing.T) {
	var argList ArgumentList

	// Empty args
	require.Empty(t, argList.Args)
	require.Equal(t, "", argList.String())

	arg1 := "my-arg1"
	require.NoError(t, argList.Set(arg1))
	require.Equal(t, 1, len(argList.Args))
	require.Equal(t, arg1, argList.String())
	require.Equal(t, arg1, argList.Args[0])

	arg2 := "another-arg"
	require.NoError(t, argList.Set(arg2))
	require.Equal(t, 2, len(argList.Args))
	require.Equal(t, arg1+arg2, argList.String())
	require.Equal(t, arg2, argList.Args[1])
}
