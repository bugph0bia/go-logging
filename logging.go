package logging

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// フォーマット指定子
const (
	FDatetime = "${Datetime}" // フォーマット指定子：日時
	FLevel    = "${Level}"    // フォーマット指定子：レベル
	FAttrs    = "${Attrs}"    // フォーマット指定子：属性リスト
	FMessage  = "${Message}"  // フォーマット指定子：メッセージ
)

// オプション
type Option struct {
	Level      slog.Level // 出力するログレベル
	WithStdout bool       // 標準出力にもログ出力するスイッチ
}

// フォーマットオプション
type Format struct {
	Line          string // ログ全体のフォーマット
	Datetime      string // 日時のフォーマット
	AttrBetween   string // 属性のキーと値の間の文字
	AttrDelimiter string // 属性と属性の間の区切り文字
	AttrPrefix    string // 属性リストの接頭辞
	AttrSuffix    string // 属性リストの接尾辞
}

// ロギングハンドラ
type Handler struct {
	slog.Handler
	RotateLogger *lumberjack.Logger
	Option       Option
	Format       Format
}

// ロギングハンドラ生成
func NewHandler(fname string) *Handler {
	// lumberjack.Logger 構造体をデフォルト値で生成
	w := &lumberjack.Logger{
		Filename:   fname,
		MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     0,
		LocalTime:  true,
		Compress:   false,
	}

	return &Handler{
		// slogのTextHandlerを生成
		Handler: slog.NewTextHandler(w, nil),

		// lubverjack.Logger 構造体も保持する
		// Close() を行う場合などに必要
		RotateLogger: w,

		// オプションのデフォルト値
		Option: Option{
			Level:      slog.LevelInfo,
			WithStdout: true,
		},

		// フォーマットのデフォルト値
		Format: Format{
			Line:          fmt.Sprintf("%s %s %s: %s", FDatetime, FLevel, FAttrs, FMessage),
			Datetime:      "2006/01/02 15:04:05",
			AttrBetween:   "=",
			AttrDelimiter: ", ",
			AttrPrefix:    "[",
			AttrSuffix:    "]",
		},
	}
}

// ロギング
// slog.Hander をベースに Handle() のみ実装する
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	// 対象レベル未満は出力しない
	if r.Level < h.Option.Level {
		return nil
	}

	// 属性を文字列に変換
	var attrs string
	as := []string{}
	r.Attrs(func(attr slog.Attr) bool {
		as = append(as, fmt.Sprintf("%s%s%v", attr.Key, h.Format.AttrBetween, attr.Value))
		return true
	})
	if len(as) > 0 {
		attrs = h.Format.AttrPrefix + strings.Join(as, h.Format.AttrDelimiter) + h.Format.AttrSuffix
	}

	// 指定されたフォーマットでログを構築
	line := h.Format.Line
	line = strings.Replace(line, FDatetime, r.Time.Format(h.Format.Datetime), -1) // 日時
	line = strings.Replace(line, FLevel, r.Level.String(), -1)                    // レベル
	line = strings.Replace(line, FAttrs, attrs, -1)                               // 属性
	line = strings.Replace(line, FMessage, r.Message, -1)                         // メッセージ
	line += "\n"

	// ログファイルへ出力
	h.RotateLogger.Write([]byte(line))

	if h.Option.WithStdout {
		// 標準出力へ出力
		fmt.Print(line)
	}

	return nil
}

// ロガー生成
func NewLoggerFromHandler(h *Handler) *slog.Logger {
	return slog.New(h)
}

// デフォルト値でロガー生成
func NewLogger(fname string) *slog.Logger {
	return slog.New(NewHandler(fname))
}
