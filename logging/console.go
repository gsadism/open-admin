package logging

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type Console struct {
	zapcore.Encoder
	Formatter func(ent Entry) (string, error)
}

func (c *Console) Clone() zapcore.Encoder {
	return &Console{
		Encoder: c.Encoder.Clone(),
	}
}

func (c *Console) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if buf, err := c.Encoder.Clone().EncodeEntry(ent, fields); err != nil {
		return nil, err
	} else {
		if str, err := c.Formatter(ent); err != nil {
			return nil, err
		} else {
			buf.Reset()
			buf.AppendString(str)
			return buf, nil
		}
	}
}
