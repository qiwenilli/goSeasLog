package goSeasLog

import (
	"testing"
)

func Benchmark_Test(b *testing.B) {

	_lint := [5]int{1, 2, 3, 4, 5}
	_lint2 := [5]string{"1", "2", "3", "4", "5"}

	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		go func() {

            _logger := New()
            //第一个参数为空，将不会把日志输出到文件
            //第二个参数是时间format，用于切分文件
            //第三个参数 true 将把不同级别的日志再分别切分
            // _logger.SetLogFile("./logs/", "200601021504", true)
            // _logger.SetLogFile("./logs/", "200601021504", false)
            _logger.SetLogFile("./logs/", "2006010215", true)
            //是否输出到终端
            _logger.SetTerminalOut(false)
            _logger.SetLevel(Debug)
            _logger.SetDateFormat("2006/1/2 15:04:05") 
            _logger.SetGap(" & ")
            //调整占位符位置
            _logger.SetLogFormat("%L %P %T %Q %H %F %M")
            _logger.SetHostName("imeiren.com")


            _logger.Debug("aabb", 1, 3, i)
            _logger.Debug("aabb", 1, 3, i, _lint, _lint2)
            _logger.Info("aabb", 1, 3, i)
            _logger.Info("aabb", 1, 3, i, _lint, _lint2)
            _logger.Warn("aabb", 1, 3, i)
            _logger.Warn("aabb", 1, 3, i, _lint, _lint2)
            _logger.Error("aabb", 1, 3, i)
            _logger.Error("aabb", 1, 3, i, _lint, _lint2)

		}()

		if i > 100000 {
			b.StopTimer()
			b.SkipNow()
		}

	}

}
