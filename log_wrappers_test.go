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
	Info("Simple info message", "key1", "value1", "keyz", "valuez")

	// Test Debug function
	Debug("Simple debug message", "key1", "value1", "key2", "value2")

	// Test Warn function
	Warn("Simple warning message", "key1", "value1")

	// Test Error function
	Error("Simple error message", "key1", "value1", "key2", "value2")
}

func TestStrArrayFromAnyArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []interface{}
		expected []string
	}{
		{
			name:     "empty args",
			args:     []interface{}{},
			expected: []string{},
		},
		{
			name:     "string arguments",
			args:     []interface{}{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "mixed type arguments",
			args:     []interface{}{123, "hello", 45.67, true, nil},
			expected: []string{"123", "hello", "45.67", "true", "<nil>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strArrayFromAnyArgs(tt.args...)

			if len(result) != len(tt.expected) {
				t.Errorf("strArrayFromAnyArgs() length = %v, expected length %v", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("strArrayFromAnyArgs()[%d] = %v, expected %v", i, v, tt.expected[i])
				}
			}
		})
	}
}
