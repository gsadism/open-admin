package logging

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type StreamHandler struct {
	zapcore.Encoder
	Formatter func(ent Entry) (string, error)
}

func (s *StreamHandler) Close() zapcore.Encoder {
	return &StreamHandler{
		Encoder: s.Encoder.Clone(),
	}
}

func (s *StreamHandler) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if buf, err := s.Encoder.Clone().EncodeEntry(ent, fields); err != nil {
		return nil, err
	} else {
		if str, err := s.Formatter(ent); err != nil {
			return nil, err
		} else {
			buf.Reset()
			buf.AppendString(str)
			return buf, nil
		}
	}
}
