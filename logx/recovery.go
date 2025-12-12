package logx

import (
	"context"
	"fmt"
	"github.com/fengzhi09/golibx/gox"
)

// Recovery will be obsolete, use RecoveryCtx, RecoveryWarn or RecoveryBase instead
func Recovery() {
	if err := recover(); err != nil {
		ctx := context.Background()
		stack := gox.Stack(1)
		logMsg := fmt.Sprintf("[Recovery] panic(%v) recovered \n trace:%s", err, stack)
		Error(ctx, logMsg)
	}
}

func RecoveryCtx(ctx context.Context) {
	if err := recover(); err != nil {
		stack := gox.Stack(1)
		logMsg := fmt.Sprintf("[RecoveryCtx] panic(%v) recovered \n trace:%s", err, stack)
		Error(ctx, logMsg)
		fmt.Println(logMsg)
	}
}

func RecoveryWarn(ctx context.Context, msg string) {
	if err := recover(); err != nil {
		stack := gox.Stack(1)
		logMsg := fmt.Sprintf("[RecoveryWarn] panic(%v) recovered \n trace:%s", msg+fmt.Sprint(err), stack)
		Warn(ctx, logMsg)
		fmt.Println(logMsg)
	}
}

func RecoveryBase(ctx context.Context, lvl LogLevel, msg string) {
	if err := recover(); err != nil {
		stack := gox.Stack(1)
		logMsg := fmt.Sprintf("[Recovery]%v panic(%v) recovered  \n trace:%s", msg, err, stack)
		switch lvl {
		case WARN:
			Warn(ctx, logMsg)
		default:
			Error(ctx, logMsg)
		}
		fmt.Println(logMsg)
	}
}
