package logging

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// ロギングハンドラ
type handler struct {
	slog.Handler
	w          io.Writer
	level      slog.Level
	withStdout bool
}

// ロギングハンドラ生成
func newHandler(w io.Writer, level slog.Level, withStdout bool) *handler {
	return &handler{
		// slog.Hander をベースに Handle() のみ実装する
		Handler:    slog.NewTextHandler(w, nil),
		w:          w,
		level:      level,
		withStdout: withStdout,
	}
}

// フォーマットしながらログ出力
func (h *handler) Handle(_ context.Context, r slog.Record) error {
	// 対象レベル未満は出力しない
	if r.Level < h.level {
		return nil
	}

	b := &bytes.Buffer{}

	// 日時
	b.WriteString(r.Time.Format("2006/01/02 15:04:05"))

	// レベル
	b.WriteString(" ")
	switch r.Level {
	case slog.LevelDebug:
		b.WriteString("DEBUG")
	case slog.LevelInfo:
		b.WriteString("INFO")
	case slog.LevelWarn:
		b.WriteString("WARN")
	case slog.LevelError:
		b.WriteString("ERROR")
	}

	// 属性
	attrs := []string{}
	r.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
		return true
	})
	if len(attrs) > 0 {
		b.WriteString(" ")
		b.WriteString("[")
		fmt.Fprint(b, strings.Join(attrs, ", "))
		b.WriteString("]")
	}

	// メッセージ
	b.WriteString(": ")
	b.WriteString(r.Message)

	// ログファイルへ出力
	b.WriteByte('\n')
	h.w.Write(b.Bytes())
	// 標準出力
	if h.withStdout {
		fmt.Print(b.String())
	}

	return nil
}

// パラメータ
var (
	// 出力するログレベル
	Level slog.Level = slog.LevelInfo

	// ログを標準出力にも出力するかどうか
	WithStdout bool = true

	// ログファイルの最大サイズ(MB)
	// このサイズを超えるとログローテーションする
	MaxSizeMB int = 1

	// バックアップファイルの最大数
	// 0 の場合は上限なし
	MaxBackups int = 10

	// バックアップファイルの最大保持日数
	// 0 の場合は上限なし
	MaxAge int = 0

	// バックアップファイルの時刻をローカルタイムにするかどうか
	LocalTime bool = true

	// バックアップファイルを gzip 圧縮するかどうか
	Compress bool = false
)

// ロガー生成
func NewLogger(fname string) *slog.Logger {
	return slog.New(newHandler(&lumberjack.Logger{
		Filename:   fname,
		MaxSize:    MaxSizeMB,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
		LocalTime:  LocalTime,
		Compress:   Compress,
	}, Level, WithStdout))
}
