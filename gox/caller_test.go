package gox

import (
	"strings"
	"testing"
)

// 不再使用全局测试上下文，直接在需要的地方创建context

func TestCallerInfo(t *testing.T) {
	// 测试CallerOfDep函数 - 这些是纯函数，不依赖于日志系统
	dep0 := CallerDep(0)
	if dep0 == nil || dep0.File == "" {
		t.Error("CallerDep(0) returned nil or invalid caller info")
	}

	dep1 := CallerDep(1)
	if dep1 == nil || dep1.File == "" {
		t.Error("CallerDep(1) returned nil or invalid caller info")
	}

	// 测试CallerOfName函数
	nameInfo := CallerLike("TestCallerInfo")
	if nameInfo == nil || nameInfo.File == "" || !strings.Contains(nameInfo.Func, "TestCallerInfo") {
		t.Error("CallerLike('TestCallerInfo') returned nil or invalid caller info")
	}

	// 测试不存在的调用者
	notFound := CallerDep(100)
	if notFound == nil || notFound != NotFound {
		t.Error("CallerDep(100) should return NotFound")
	}
}

func TestStack(t *testing.T) {
	// 测试Stack函数 - 这是纯函数，不依赖于日志系统
	stack := Stack(0)
	if len(stack) == 0 {
		t.Error("Stack(0) returned empty stack")
	}

	// 验证堆栈包含当前函数
	stackStr := string(stack)
	if !strings.Contains(stackStr, "TestStack") {
		t.Error("Stack does not contain current function name")
	}

	// 测试跳过堆栈帧
	stackSkip := Stack(1)
	stackSkipStr := string(stackSkip)
	if stackSkipStr == stackStr {
		t.Error("Stack(1) should skip one frame but didn't")
	}
}
