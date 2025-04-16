package email_sender

import "io"

// lineBreakByteLimit defines the maximum number of bytes per line
// before inserting a CRLF line break (used for base64-encoded MIME content).
const lineBreakByteLimit = 76

// Base64LineWriter is a writer that wraps another io.Writer and inserts
// CRLF ("\r\n") line breaks every "lineBreakByteLimit" bytes, which is a standard requirement
// for base64-encoded content in MIME messages.
type Base64LineWriter struct {
	w      io.Writer
	buffer []byte
	count  int
}

// NewBase64LineWriter creates a new Base64LineWriter that writes to the given
// io.Writer. It ensures that lines do not exceed "lineBreakByteLimit" bytes by inserting
// a CRLF after each full chunk.
func NewBase64LineWriter(w io.Writer) io.Writer {
	return &Base64LineWriter{
		w:      w,
		buffer: make([]byte, 0, lineBreakByteLimit),
	}
}

// Write writes the provided bytes to the underlying writer,
// automatically inserting CRLF line breaks every "lineBreakByteLimit" bytes
// to conform to base64 MIME encoding rules.
func (lw *Base64LineWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		lw.buffer = append(lw.buffer, b)
		lw.count++
		if lw.count == lineBreakByteLimit {
			if _, err := lw.w.Write(lw.buffer); err != nil {
				return 0, err
			}
			if _, err := lw.w.Write([]byte("\r\n")); err != nil {
				return 0, err
			}
			lw.buffer = lw.buffer[:0]
			lw.count = 0
		}
	}
	return len(p), nil
}

// Close flushes any remaining bytes in the buffer and appends a final CRLF.
// This should be called when all data has been written to ensure complete output.
func (lw *Base64LineWriter) Close() error {
	if len(lw.buffer) > 0 {
		if _, err := lw.w.Write(lw.buffer); err != nil {
			return err
		}
		if _, err := lw.w.Write([]byte("\r\n")); err != nil {
			return err
		}
	}
	return nil
}
