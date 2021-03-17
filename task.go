package gorun

import (
	"context"
	"time"
)

type Task interface {
	NextTime(now time.Time) time.Time // 次回実行時刻
	Run(ctx context.Context)          // 実行される処理
}
