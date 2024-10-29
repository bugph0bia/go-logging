package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// テスト用のログファイル名
const logfile = "test.log"

// 最終ログを取得
func lastLog() string {
	f, err := os.Open(logfile)
	if err != nil {
		return ""
	}
	defer f.Close()

	// 最終行まで読み捨てる
	// テスト上はログファイルは大きくならないためこの方法で問題なし
	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	return line
}

// 標準出力のキャプチャを開始
func startStdoutCapture() (r *os.File, w *os.File, orgStdout *os.File) {
	// パイプを作成して標準出力をフック
	r, w, _ = os.Pipe()
	orgStdout = os.Stdout
	os.Stdout = w
	return
}

// 標準出力のキャプチャを終了
func endStdoutCapture(r *os.File, w *os.File, orgStdout *os.File) string {
	// 標準出力のフックを終了
	os.Stdout = orgStdout
	w.Close()

	// パイプを通して標準出力の内容を取得する
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// デフォルトロガーのテスト
// NewLogger()を呼び出した場合も同じロジックのため、そちらのテストも兼ねる
func TestDefault(t *testing.T) {
	handler := NewHandler(logfile)
	t.Cleanup(func() {
		// lumberjack がファイルを開いているため Close してから削除する必要がある
		handler.RotateLogger.Close()
		os.Remove(logfile)
	})
	logger := NewLoggerFromHandler(handler)

	// 標準出力フック開始
	r, w, orgStdout := startStdoutCapture()

	// Debeg は出力されないことの確認
	logger.Debug("message")
	wantDebug := ""
	assert.Equal(t, wantDebug, lastLog())

	// Info は出力されることの確認
	// 属性1つのテスト
	logger.Info("message", "attr1", 1)
	wantInfo := fmt.Sprintf("%s INFO [attr1=1]: message", time.Now().Format("2006/01/02 15:04:05"))
	assert.Equal(t, wantInfo, lastLog())

	// Warn は出力されることの確認
	// 属性2つのテスト
	logger.Warn("message", "attr1", 1, "attr2", "v")
	wantWarn := fmt.Sprintf("%s WARN [attr1=1, attr2=v]: message", time.Now().Format("2006/01/02 15:04:05"))
	assert.Equal(t, wantWarn, lastLog())

	// Error は出力されることの確認
	// 属性なしのテスト
	logger.Error("message")
	wantError := fmt.Sprintf("%s ERROR : message", time.Now().Format("2006/01/02 15:04:05"))
	assert.Equal(t, wantError, lastLog())

	// 標準出力フック終了
	actualStdout := endStdoutCapture(r, w, orgStdout)

	// 標準出力にログが出力されることの確認
	wantStdout := wantInfo + "\n" + wantWarn + "\n" + wantError + "\n"
	assert.Equal(t, wantStdout, actualStdout)
}

// Option.Level のテスト
func TestOptionLevel(t *testing.T) {
	handler := NewHandler(logfile)
	t.Cleanup(func() {
		// lumberjack がファイルを開いているため Close してから削除する必要がある
		handler.RotateLogger.Close()
		os.Remove(logfile)
	})
	handler.Option.Level = slog.LevelError // デフォルトとは異なる値でテスト
	logger := NewLoggerFromHandler(handler)

	// 標準出力フック開始
	r, w, orgStdout := startStdoutCapture()

	// Debug, Info, Warn は出力されないことの確認
	logger.Debug("message")
	wantDebug := ""
	assert.Equal(t, wantDebug, lastLog())

	logger.Info("message")
	wantInfo := ""
	assert.Equal(t, wantInfo, lastLog())

	logger.Warn("message")
	wantWarn := ""
	assert.Equal(t, wantWarn, lastLog())

	// Error は出力されることの確認
	logger.Error("message")
	wantError := fmt.Sprintf("%s ERROR : message", time.Now().Format("2006/01/02 15:04:05"))
	assert.Equal(t, wantError, lastLog())

	// 標準出力フック終了
	endStdoutCapture(r, w, orgStdout)
}

// Option.WithStdout のテスト
// NewLogger()を呼び出した場合も同じロジックのため、そちらのテストも兼ねる
func TestOptionWithStdout(t *testing.T) {
	handler := NewHandler(logfile)
	t.Cleanup(func() {
		// lumberjack がファイルを開いているため Close してから削除する必要がある
		handler.RotateLogger.Close()
		os.Remove(logfile)
	})
	handler.Option.WithStdout = false // デフォルトとは異なる値でテスト
	logger := NewLoggerFromHandler(handler)

	// 標準出力フック開始
	r, w, orgStdout := startStdoutCapture()

	// Debeg は出力されない
	logger.Debug("message")
	// Info, Warn, Error は出力される
	logger.Info("message", "attr1", 1)
	logger.Warn("message", "attr1", 1, "attr2", "v")
	logger.Error("message")

	// 標準出力フック終了
	actualStdout := endStdoutCapture(r, w, orgStdout)

	// 標準出力にログが出力されないことの確認
	wantStdout := ""
	assert.Equal(t, wantStdout, actualStdout)
}

// Format のテスト
// NewLogger()を呼び出した場合も同じロジックのため、そちらのテストも兼ねる
func TestFormat(t *testing.T) {
	handler := NewHandler(logfile)
	t.Cleanup(func() {
		// lumberjack がファイルを開いているため Close してから削除する必要がある
		handler.RotateLogger.Close()
		os.Remove(logfile)
	})
	// デフォルトとは異なる値でテスト
	handler.Format.Line = fmt.Sprintf("%s %s %s: %s", FMessage, FAttrs, FLevel, FDatetime) // 順番を逆転
	handler.Format.Datetime = "06-01-02 15_04_05"
	handler.Format.AttrBetween = "=>"
	handler.Format.AttrDelimiter = " | "
	handler.Format.AttrPrefix = "{ "
	handler.Format.AttrSuffix = " }"
	logger := NewLoggerFromHandler(handler)

	// 標準出力フック開始
	r, w, orgStdout := startStdoutCapture()

	// Warn を例に指定フォーマットで出力されることを確認
	logger.Warn("message")
	wantWarn := fmt.Sprintf("message  WARN: %s", time.Now().Format("06-01-02 15_04_05"))
	assert.Equal(t, wantWarn, lastLog())

	// Error を例に指定フォーマットで出力されることを確認
	logger.Error("message", "attr1", 1, "attr2", "v")
	wantError := fmt.Sprintf("message { attr1=>1 | attr2=>v } ERROR: %s", time.Now().Format("06-01-02 15_04_05"))
	assert.Equal(t, wantError, lastLog())

	// 標準出力フック終了
	endStdoutCapture(r, w, orgStdout)
}

// lumberjack.Logger のパラメータのテストは省略
