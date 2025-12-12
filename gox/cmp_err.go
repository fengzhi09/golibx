package gox

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrDiffCore = errors.New("diff")

type ErrDiff struct {
	core error
	msg  string
}

func NewErrDiff(msg string) error {
	return &ErrDiff{core: ErrDiffCore, msg: msg}
}

func NewErrDiff2(format string, args ...any) error {
	return &ErrDiff{core: ErrDiffCore, msg: fmt.Sprintf(format, args...)}
}

func (e *ErrDiff) Error() string {
	return e.msg
}

func (e *ErrDiff) Is(err error) bool {
	return e.core == err
}
