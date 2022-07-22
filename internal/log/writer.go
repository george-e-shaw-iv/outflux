package log

import (
	"io"
	"os"
	"sync"
)

// writers are the writers used by the logger. These can be defined per level
// and if a writer isn't defined for a level then it falls back to the default
// writer.
type writers struct {
	DefaultWriter io.Writer
	Levels        map[int]io.Writer
}

// defaultWriters returns the default set of writers used by the Logger.
func defaultWriters() *writers {
	stdout := newSyncWriter(os.Stdout)
	stderr := newSyncWriter(os.Stderr)

	return &writers{
		DefaultWriter: stdout,
		Levels: map[int]io.Writer{
			LevelDebug: stdout,
			LevelInfo:  stdout,
			LevelWarn:  stdout,
			LevelError: stderr,
			LevelFatal: stderr,
		},
	}
}

// SetWriterOnAll overrides all writers on each level and the default writer to the
// provided writer.
func (w *writers) SetWriterOnAll(writer io.Writer) {
	sWriter := newSyncWriter(writer)

	w.DefaultWriter = sWriter
	for k := range w.Levels {
		w.Levels[k] = sWriter
	}
}

// SetWriterOnLevels overrides the writer on each provided level to the writer that
// is provided.
func (w *writers) SetWriterOnLevels(writer io.Writer, levels ...int) {
	sWriter := newSyncWriter(writer)

	for _, level := range levels {
		if _, ok := w.Levels[level]; ok {
			w.Levels[level] = sWriter
		}
	}
}

// ByLevel retrieves a writer for the specified level, falling back to the default
// writer if a writer does not exist at that level.
func (w *writers) ByLevel(level int) io.Writer {
	writer := w.Levels[level]
	if writer == nil {
		writer = w.DefaultWriter
	}
	return writer
}

// syncWriter protects against fragmented log lines.
type syncWriter struct {
	mu     sync.Mutex
	writer io.Writer
}

// newSyncWriter returns a syncWriter using the passed in io.Writer.
func newSyncWriter(writer io.Writer) *syncWriter {
	return &syncWriter{
		writer: writer,
	}
}

// Writer implements the io.Writer interface for syncWriter.
func (s *syncWriter) Write(b []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.writer.Write(b)
}
