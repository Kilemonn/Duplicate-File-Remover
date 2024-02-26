package argument_list

import (
	"strings"
)

type ArgumentList struct {
	Args []string
}

func (list *ArgumentList) String() string {
	return strings.Join(list.Args, "")
}

func (list *ArgumentList) Set(value string) error {
	list.Args = append(list.Args, value)
	return nil
}
