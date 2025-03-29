package cmd

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput bool
		wantError  bool
	}{
		{
			name:       "GIVEN no arguments THEN help is displayed without error",
			args:       []string{},
			wantOutput: true,
			wantError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()

			if (err != nil) != tc.wantError {
				t.Errorf("Execute() error = %v, wantError %v", err, tc.wantError)
			}

			gotOutput := buf.String()
			if tc.wantOutput && len(gotOutput) == 0 {
				t.Error("Expected non-empty output but got nothing")
			}

			if !tc.wantOutput && len(gotOutput) > 0 {
				t.Errorf("Expected no output but got: %s", gotOutput)
			}
		})
	}
}

func TestCommandRegistration(t *testing.T) {
	tests := []struct {
		name       string
		commandUse string
		wantFound  bool
	}{
		{
			name:       "GIVEN add command THEN it is registered in root command",
			commandUse: "add",
			wantFound:  true,
		},
		{
			name:       "GIVEN commit command THEN it is registered in root command",
			commandUse: "commit",
			wantFound:  true,
		},
		{
			name:       "GIVEN nonexistent command THEN it is not found in root command",
			commandUse: "nonexistent",
			wantFound:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			found := false
			for _, cmd := range rootCmd.Commands() {
				if cmd.Use == tc.commandUse {
					found = true
					break
				}
			}

			if found != tc.wantFound {
				t.Errorf("Command %q registration = %v, want %v", tc.commandUse, found, tc.wantFound)
			}
		})
	}
}
