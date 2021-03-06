package dsky

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	fc "github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/gosuri/uitable/util/strutil"
)

var (
	MaxLogLabelWidth = 7
	MaxLogLineWidth  = 140
)

type interactiveLogger struct {
	module string
	out    io.Writer
	LogAction
}

func NewInteractiveLogger(out io.Writer) *interactiveLogger {
	return &interactiveLogger{out: out}
}

func (l *interactiveLogger) WithModule(module string) Logger {
	l.module = module
	return l
}

func (l *interactiveLogger) WithAction(action LogAction) Logger {
	l.LogAction = action
	return l
}

func (l *interactiveLogger) Info(msg ...interface{}) LogItem {
	lm := &logItem{
		logModeType: logModeTypeInfo,
		LogAction:   l.LogAction,
		labelColor:  Color.Success,
		msgColor:    Color.Hi,
		msg:         msg,
		module:      l.module,
	}
	l.writelog(lm)
	return lm
}

func (l *interactiveLogger) Debug(msg ...interface{}) LogItem {
	lm := &logItem{
		logModeType: logModeTypeDebug,
		LogAction:   l.LogAction,
		labelColor:  Color.Hi,
		msgColor:    Color.Hi,
		msg:         msg,
		module:      l.module,
	}
	l.writelog(lm)
	return lm
}

func (l *interactiveLogger) Warn(msg ...interface{}) LogItem {
	lm := &logItem{
		logModeType: logModeTypeWarn,
		LogAction:   l.LogAction,
		labelColor:  Color.Notice,
		msgColor:    Color.Normal,
		msg:         msg,
		module:      l.module,
	}
	l.writelog(lm)
	return lm
}

func (l *interactiveLogger) Error(msg ...interface{}) LogItem {
	lm := &logItem{
		logModeType: logModeTypeError,
		LogAction:   l.LogAction,
		labelColor:  Color.Failure,
		msgColor:    Color.Normal,
		msg:         msg,
		module:      l.module,
	}
	l.writelog(lm)
	return lm
}

func (l *interactiveLogger) writelog(lm LogItem) {
	b, _ := lm.String()
	fmt.Fprintln(l.out, b)
}

type logItem struct {
	logModeType
	LogAction
	labelColor *fc.Color
	msgColor   *fc.Color
	msg        []interface{}
	module     string
}

func (lm *logItem) WithLabelColor(color *fc.Color) *logItem {
	lm.labelColor = color
	return lm
}

func (lm *logItem) WithMessageColor(color *fc.Color) *logItem {
	lm.msgColor = color
	return lm
}

func (lm *logItem) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	label := lm.labelColor.Sprintf("(%s)", lm.logModeType)
	if len(lm.LogAction) > 0 {
		label = lm.labelColor.Sprintf("(%s)", lm.LogAction)
	}
	label = strutil.PadRight(label, MaxLogLabelWidth, ' ')
	msg := []string{}

	if len(lm.module) > 0 {
		if lm.msgColor != nil {
			msg = append(msg, lm.msgColor.Sprintf("[%s]", lm.module))
		} else {
			msg = append(msg, fmt.Sprintf("[%s]", lm.module))
		}
	}

	for _, m := range lm.msg {
		if lm.msgColor != nil {
			msg = append(msg, lm.msgColor.Sprintf("%v", m))
		} else {
			msg = append(msg, fmt.Sprintf("%v", m))
		}
	}

	buf.WriteString(prefixedMsg(label, strings.Join(msg, " "), MaxLogLineWidth))
	return buf.Bytes(), nil
}

func (lm *logItem) String() (string, error) {
	b, err := lm.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func prefixedMsg(label string, msg string, width int) string {
	var buf bytes.Buffer
	cell := &uitable.Cell{Width: uint(width), Wrap: true, Data: msg}

	var lines []string
	for i, line := range strings.Split(cell.String(), "\n") {
		var lb bytes.Buffer
		if i == 0 {
			lb.WriteString(label)
			lb.WriteString(" ")
			lb.WriteString(line)
		} else {
			s := strutil.PadLeft(line, strutil.StringWidth(label)+strutil.StringWidth(line)+1, ' ')
			lb.WriteString(s)
		}
		lines = append(lines, lb.String())
	}
	buf.WriteString(strings.Join(lines, "\n"))
	return buf.String()
}

func prefixedMsgT(label string, msg string, width uint) string {
	t := uitable.New().AddRow(label, msg)
	t.MaxColWidth = width
	t.Wrap = true
	return t.String()
}
