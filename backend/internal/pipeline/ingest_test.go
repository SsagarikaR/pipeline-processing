package pipeline

import (
	"testing"
)

func TestOpenSource_PathTraversalProtection(t *testing.T) {
	// Restore SandboxDir at the end of the test
	originalSandbox := SandboxDir
	defer func() { SandboxDir = originalSandbox }()

	// Set a mock sandbox directory
	SandboxDir = "/data/inputs"

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid path inside sandbox",
			path:        "/data/inputs/my_file.csv",
			expectError: true, // It will error because file doesn't actually exist, but NOT a traversal error
			errorMsg:    "no such file or directory",
		},
		{
			name:        "Path traversal with .. inside sandbox",
			path:        "/data/inputs/../outputs/file.csv",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
		{
			name:        "Absolute path outside sandbox",
			path:        "/etc/passwd",
			expectError: true,
			errorMsg:    "path outside allowed sandbox /data/inputs",
		},
		{
			name:        "Relative path climbing",
			path:        "../../../etc/passwd",
			expectError: true,
			errorMsg:    "path traversal detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, closeFn, err := openSource(tt.path)
			
			if closeFn != nil {
				defer closeFn()
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else {
					// We just want to check if the error contains the expected message,
					// or for the "no such file" we just ensure it didn't fail the traversal check.
					// Since os.Open on a non-existent file returns a PathError, we check the suffix.
					if !containsErrStr(err.Error(), tt.errorMsg) {
						t.Errorf("expected error to contain %q, got: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			}
		})
	}
}

func containsErrStr(actual, expected string) bool {
	// Simple helper to check if expected substring is in actual error string
	return len(actual) >= len(expected) && (actual[len(actual)-len(expected):] == expected || contains(actual, expected))
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
