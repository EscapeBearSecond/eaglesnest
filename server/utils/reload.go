package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func Reload() error {
	if runtime.GOOS == "windows" {
		return errors.New("系统不支持")
	}
	pid := os.Getpid()
	cmd := exec.Command("kill", "-1", strconv.Itoa(pid))
	return cmd.Run()
}

func RemoveValueFromList(rdb redis.UniversalClient, listName string, targetValue string) bool {
	var ctx = context.Background()
	for {
		// 阻塞获取队列中的值
		result, err := rdb.LPop(ctx, listName).Result()
		if err != nil {
			fmt.Println("取出值失败:", err)
			return false
		}
		value := result
		if value == targetValue {
			// 找到目标值，成功取出
			return true
		} else {
			// 不是目标值，重新放回队列的尾部
			rdb.RPush(ctx, listName, value)
		}
	}
}
