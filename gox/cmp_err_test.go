package gox

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrDiff(t *testing.T) {
	// 测试基本错误创建
	err := NewErrDiff("test difference")
	assert.NotNil(t, err)
	assert.Equal(t, "test difference", err.Error())

	// 测试错误类型断言
	errDiff, ok := err.(*ErrDiff)
	assert.True(t, ok)
	assert.Equal(t, ErrDiffCore, errDiff.core)
	assert.Equal(t, "test difference", errDiff.msg)

	// 测试空消息
	emptyErr := NewErrDiff("")
	assert.Equal(t, "", emptyErr.Error())
}

func TestNewErrDiff2(t *testing.T) {
	// 测试带格式的错误创建
	err := NewErrDiff2("diff at index %d: %s vs %s", 5, "actual", "expected")
	assert.NotNil(t, err)
	assert.Equal(t, "diff at index 5: actual vs expected", err.Error())

	// 测试错误类型断言
	errDiff, ok := err.(*ErrDiff)
	assert.True(t, ok)
	assert.Equal(t, ErrDiffCore, errDiff.core)
	assert.Equal(t, "diff at index 5: actual vs expected", errDiff.msg)

	// 测试无参数格式
	noArgsErr := NewErrDiff2("simple message")
	assert.Equal(t, "simple message", noArgsErr.Error())
}

func TestErrDiff_Is(t *testing.T) {
	// 测试Is方法
	err := NewErrDiff("test error")
	errDiff, _ := err.(*ErrDiff)

	// 与ErrDiffCore比较
	assert.True(t, errDiff.Is(ErrDiffCore))
	assert.True(t, errors.Is(err, ErrDiffCore)) // 测试标准库的errors.Is

	// 与其他错误比较
	otherErr := errors.New("other error")
	assert.False(t, errDiff.Is(otherErr))
	assert.False(t, errors.Is(err, otherErr))

	// 与nil比较
	assert.False(t, errDiff.Is(nil))
	assert.False(t, errors.Is(err, nil))
}

func TestErrDiff_Error(t *testing.T) {
	// 测试各种消息的Error方法
	cases := []struct {
		msg      string
		expected string
	}{
		{"simple message", "simple message"},
		{"", ""},
		{"12345", "12345"},
		{"特殊字符: !@#$%^&*()", "特殊字符: !@#$%^&*()"},
	}

	for _, tc := range cases {
		err := NewErrDiff(tc.msg)
		assert.Equal(t, tc.expected, err.Error())
	}
}
