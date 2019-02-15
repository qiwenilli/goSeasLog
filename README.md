# goSeasLog

是一个模拟php Seaslog 的一个go实现日志输出库

```
	_logger := goSeasLog.New()
	
	//第一个参数为空，将不会把日志输出到文件
	//第二个参数是时间format，用于切分文件
	//第三个参数 true 将把不同级别的日志再分别切分
    // _logger.SetLogFile("./logs/", "200601021504", true)
    // _logger.SetLogFile("./logs/", "200601021504", false)
    _logger.SetLogFile("./logs/", "2006010215", true)
	_logger.SetLogFile("./logs/", "2006010215", true)
	
	//是否输出到终端
	_logger.SetTerminalOut(false)
	_logger.SetLevel(goSeasLog.Debug)
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

```


