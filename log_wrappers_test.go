package logger

import "testing"

func TestFormattingFunctions(t *testing.T) {
	InitLog(LogConfig{
		Formatter: "json",
		LogLevel:  "debug",
	})
	defer CloseLog()

	// Test F function
	F("Number: %d String: %s", 42, "test")

	// Test InfoF function
	InfoF("Info message with %s and %d", "string", 123)

	// Test DebugF function
	DebugF("Debug message with %.2f", 3.14159)

	// Test WarnF function
	WarnF("Warning: %s occurred %d times", "event", 5)

	// Test ErrorF function
	ErrorF("Error code: %x", 255)
}

func TestSimpleLogFunctions(t *testing.T) {
	InitLog(LogConfig{
		Formatter: "json",
		LogLevel:  "debug",
	})
	defer CloseLog()

	// Test Info function
	Info("Simple info message", "key1", "value1")

	// Test Debug function
	Debug("Simple debug message", "key1", "value1", "key2", "value2")

	// Test Warn function
	Warn("Simple warning message", "key1", "value1")

	// Test Error function
	Error("Simple error message", "key1", "value1", "key2", "value2")
}

func TestAttrLogFunctions(t *testing.T) {
	InitLog(LogConfig{
		Formatter: "json",
		LogLevel:  "debug",
	})
	defer CloseLog()

	a := Attr{"key1", "value1"}
	InfoAttrs("A log", a, Attr{"key2", "value2"}, Attr{"key3", 1})

	// Test DebugAttrs
	DebugAttrs("Debug with attributes",
		Attr{"string_key", "string_value"},
		Attr{"int_key", 42},
		Attr{"float_key", 3.14})

	// Test WarnAttrs
	WarnAttrs("Warning with attributes",
		Attr{"warning_code", 301},
		Attr{"source", "test"})

	// Test ErrorAttrs
	ErrorAttrs("Error with attributes",
		Attr{"error_code", 500},
		Attr{"system", "database"},
		Attr{"retry", false})
}
