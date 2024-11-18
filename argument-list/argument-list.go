package argument_list

import (
	"strings"
)

// Satisfies the [flag.Value] interface to be used in flag.Var() to collect argument parameters.
// Multiple arguments that are provided with the same flag will be appended together into this struct.
type ArgumentList struct {
	Args []string
}

// String implementing for [flag.Value.String], returns all of the elements joined into a single string.
func (list *ArgumentList) String() string {
	return strings.Join(list.Args, "")
}

// Set implementing for [flag.Value.Set], will append the provided argument to the internal array.
func (list *ArgumentList) Set(value string) error {
	list.Args = append(list.Args, value)
	return nil
}
