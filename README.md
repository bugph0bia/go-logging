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

オプションを変更する場合は、まず `Logging.NewHandler` でハンドラを取得して変更してから `logging.NewLoggerFromHandler` を呼び出してロガーを生成する。  
下記サンプルコードは変更可能なオプションとデフォルト値の例。  

```go
func main() {
    // ハンドラ生成
    handler := logging.NewHandler("log.txt")

    // オプションを変更

    // 出力するログレベル
    handler.Option.Level = slog.LevelInfo
    // ログを標準出力にも出力するかどうか
    handler.WithStdout = true

    // フォーマットを変更

    // ログ全体のフォーマット
    handler.Format.Line = fmt.Sprintf("${Datetime} ${Level} ${Attrs}: ${Message}")
    // ログ全体のフォーマット（定数を利用）
    handler.Format.Line = fmt.Sprintf("%s %s %s: %s", logging.FDatetime, logging.FLevel, logging.FAttrs, logging.FMessage)
    // 日時のフォーマット
    handler.Format.Datetime = "2006/01/02 15:04:05"
    // 属性のキーと値の間の文字
    handler.Format.AttrBetween = "="
    // 属性と属性の間の区切り文字
    handler.Format.AttrDelimiter = ", "
    // 属性リストの接頭辞
    handler.Format.AttrPrefix = "["
    // 属性リストの接尾辞
    handler.Format.AttrSufix = "]"

    // ローテーションの動作を変更

    // ログファイルの最大サイズ(MB)
    // このサイズを超えるとログローテーションする
    handler.RotateLogger.MaxSizeMB = 1
    // バックアップファイルの最大数
    // 0 の場合は上限なし
    handler.RotateLogger.MaxBackups = 10
    // バックアップファイルの最大保持日数
    // 0 の場合は上限なし
    handler.RotateLogger.MaxAge = 0
    // バックアップファイルの時刻をローカルタイムにするかどうか
    handler.RotateLogger.LocalTime = true
    // バックアップファイルを gzip 圧縮するかどうか
    handler.RotateLogger.Compress = false

    // パラメータ設定したハンドラーからロガー生成
    logger := logging.NewLoggerFromHandler(handler)
}
```