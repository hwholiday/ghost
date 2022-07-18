

###   

#### 配置

- PrintTime TimeFormat // 打印的时间格式
- LogFileDir string // 日志路径
- AppName string // APP名字
- MaxSize int //文件多大开始切分
- MaxBackups int //保留文件个数
- MaxAge int //文件保留最大实际
- Level string // 日志打印等级
- CtxKey string //通过 ctx 传递 hlog 信息
- WriteFile bool // 是否写入文件
- WriteConsole bool // 是否控制台打印

##### 实现自动切割文件核心代码

```base
 zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.opts.LogFileDir + "/" + l.opts.AppName + ".log",
		MaxSize:    l.opts.MaxSize,
		MaxBackups: l.opts.MaxBackups,
		MaxAge:     l.opts.MaxAge,
		Compress:   true,
		LocalTime:  true,
	})
```

#### 实现可传递 trace 信息核心代码

```base
func (l *Logger) GetCtx(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.opts.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.Logger
}

func (l *Logger) AddCtx(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
	log := l.With(field...)
	ctx = context.WithValue(ctx, l.opts.CtxKey, log)
	return ctx, log
}
```

#### 例子（普通）

```base
hlog.NewLogger()
hlog.GetLogger().Info("hconf example success")

{"level":"info","ts":"2022-07-18 23:29:02","caller":"hlog/zap.go:36","msg":"[initLogger] zap plugin initializing completed"}
{"level":"info","ts":"2022-07-18 23:29:02","caller":"hlog/zap_test.go:12","msg":"hconf example success"}

```

#### 例子（gin）

```base
func AddTraceId() gin.HandlerFunc {
	return func(g *gin.Context) {
		traceId := g.GetHeader("traceId")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		ctx, log := hlog.GetLogger().AddCtx(g.Request.Context(), zap.Any("traceId", traceId))
		g.Request = g.Request.WithContext(ctx)
		log.Info("AddTraceId success")
		g.Next()
	}
}

log := hlog.GetLogger().GetCtx(context.Request.Context())
		log.Info("test")
		log.Debug("test")	
```

#### 例子（gin）开发模式

```base	
hlog.NewLogger()	

curl http://127.0.0.1:8888/test

{"level":"info","ts":"2022-07-18 23:48:16","caller":"example/main.go:35","msg":"hconf example success"}
{"level":"info","ts":"2022-07-18 23:48:21","caller":"example/main.go:19","msg":"AddTraceId success","traceId":"c7cd8bed-10d8-449b-ab17-338ff2651463"}
{"level":"info","ts":"2022-07-18 23:48:21","caller":"example/main.go:31","msg":"test","traceId":"c7cd8bed-10d8-449b-ab17-338ff2651463"}
{"level":"debug","ts":"2022-07-18 23:48:21","caller":"example/main.go:32","msg":"test","traceId":"c7cd8bed-10d8-449b-ab17-338ff2651463"}

```

#### 例子（gin）生产模式

```base
hlog.NewLogger(hlog.SetPrintTime(hlog.PrintTimestamp))

curl http://127.0.0.1:8888/test

{"level":"info","ts":1658159389991,"caller":"example/main.go:35","msg":"hconf example success"}
{"level":"info","ts":1658159391161,"caller":"example/main.go:19","msg":"AddTraceId success","traceId":"c304cdb9-7d54-4d79-9c33-691cfc868d3c"}
{"level":"info","ts":1658159391161,"caller":"example/main.go:31","msg":"test","traceId":"c304cdb9-7d54-4d79-9c33-691cfc868d3c"}
{"level":"debug","ts":1658159391161,"caller":"example/main.go:32","msg":"test","traceId":"c304cdb9-7d54-4d79-9c33-691cfc868d3c"}
```
