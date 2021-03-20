package gorun

import (
	"context"
	"time"
)

type Task interface {
	NextTime(now time.Time) time.Duration // 次回実行時刻までの間隔
	Run(ctx context.Context)              // 実行される処理
}
