# go-logging

slog と lumberjack を組み合わせた簡易ロギング＆ログローテーション  


## 使用方法

```go
package main

import "github.com/bugph0bia/go-logging"

func main() {
    // デフォルトパラメータでロガー生成
    logger := logging.NewLogger("log.txt")
    // ログ出力
    logger.Debug("message")
    logger.Info("message", "attr1", 10)
    logger.Warn("message", "attr1", 10, "attr2", 20)
    logger.Error("message")
}

// 以下のフォーマットでログが出力される
//
// 2024/10/24 11:22:33 DEBUG: message
// 2024/10/24 11:22:33 INFO [attr1=10]: message
// 2024/10/24 11:22:33 WARN [attr1=10, attr2=20]: message
// 2024/10/24 11:22:33 ERROR: message
```

パラメータを変更する場合は `logging.NewLogger` を呼び出す前に値を変更しておく。  
下記は設定できるパラメータとデフォルト値の例。  

```go
// 出力するログレベル
logging.Level = slog.LevelInfo

// ログを標準出力にも出力するかどうか
logging.WithStdout = true

// ログファイルの最大サイズ(MB)
// このサイズを超えるとログローテーションする
logging.MaxSizeMB = 1

// バックアップファイルの最大数
// 0 の場合は上限なし
logging.MaxBackups = 10

// バックアップファイルの最大保持日数
// 0 の場合は上限なし
logging.MaxAge = 0

// バックアップファイルの時刻をローカルタイムにするかどうか
logging.LocalTime = true

// バックアップファイルを gzip 圧縮するかどうか
logging.Compress = false
```