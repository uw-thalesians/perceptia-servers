package utility

import (
	"context"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-redis/redis"
)

func PingRedis(ctx context.Context, rc *redis.Client, sleepFailTime time.Duration,
	sleepTestTime time.Duration, logger kitlog.Logger,
	statusNotOkay chan bool) {
	for {
		select {
		case <-ctx.Done():
			_ = logger.Log("PingRedis", "ping check canceled")
			return
		default:
			if res := rc.Ping(); res.Err() != nil {
				_ = logger.Log("func", "utility.PingRedis", "pingError", res.Err(), "note", "will retry in "+sleepTestTime.String())
				select {
				case _, _ = <-statusNotOkay:
					break
				default:
					break
				}
				statusNotOkay <- true
				time.Sleep(sleepFailTime)
			} else {
				select {
				case _, _ = <-statusNotOkay:
					break
				default:
					break
				}
				statusNotOkay <- false
				time.Sleep(sleepTestTime)
			}
		}
	}
}
