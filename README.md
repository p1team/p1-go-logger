# p1-go-logger

Go標準の`log/slog`を利用した共通ロギング部品です。以下の出力先にリクエストID付きのログを出力します。

- STDOUT
- File

## 目的

ログ出力方法の共通化

## インストール

```bash
go get github.com/p1team/p1-go-logger
```

## 使い方

ログの出力先はLogger生成時に指定します。

```go
requestID := uuid.NewString()

logger := logging.New(
	logging.StdoutHandler(slog.LevelInfo),
	requestID,
)

run(ctx, logger)
```

Loggerは、関数には引数として渡し、structにはコンストラクタで渡します。

関数の場合は下記のようにします。

```go
func run(ctx context.Context, logger *slog.Logger) int {
	logger.InfoContext(ctx, "application started")
	return 0
}
```

structの場合は下記のようにします。

```go
type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) Execute(ctx context.Context) {
	s.logger.InfoContext(ctx, "execute")
}
```

## サンプルコード

### STDOUT

```go
logger := logging.New(
	logging.StdoutHandler(slog.LevelInfo),
	requestID,
)

run(ctx, logger)
```

### File

```go
handler, closer, err := logging.FileHandler(
	"./app.log",
	slog.LevelInfo,
)
if err != nil {
	return err
}
defer closer.Close()

logger := logging.New(
	handler,
	requestID,
)

run(ctx, logger)
```
