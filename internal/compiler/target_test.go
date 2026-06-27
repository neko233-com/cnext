package compiler

import (
	"testing"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantArch string
		wantOS   string
		wantEnv  string
		wantErr  bool
	}{
		{
			name:     "common linux-x64",
			input:    "linux-x64",
			wantArch: "x86_64",
			wantOS:   "linux",
			wantEnv:  "gnu",
		},
		{
			name:     "common linux-arm64",
			input:    "linux-arm64",
			wantArch: "aarch64",
			wantOS:   "linux",
			wantEnv:  "gnu",
		},
		{
			name:     "common windows-x64",
			input:    "windows-x64",
			wantArch: "x86_64",
			wantOS:   "windows",
			wantEnv:  "gnu",
		},
		{
			name:     "common macos-arm64",
			input:    "macos-arm64",
			wantArch: "aarch64",
			wantOS:   "darwin",
			wantEnv:  "none",
		},
		{
			name:     "full triple",
			input:    "x86_64-unknown-linux-gnu",
			wantArch: "x86_64",
			wantOS:   "linux",
			wantEnv:  "gnu",
		},
		{
			name:     "aarch64 triple",
			input:    "aarch64-unknown-linux-musl",
			wantArch: "aarch64",
			wantOS:   "linux",
			wantEnv:  "musl",
		},
		{
			name:     "three parts",
			input:    "x86_64-linux-musl",
			wantArch: "x86_64",
			wantOS:   "linux",
			wantEnv:  "musl",
		},
		{
			name:     "two parts",
			input:    "aarch64-darwin",
			wantArch: "aarch64",
			wantOS:   "darwin",
			wantEnv:  "gnu",
		},
		{
			name:    "single part - error",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty - error",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := ParseTarget(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseTarget(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseTarget(%q) unexpected error: %v", tt.input, err)
			}
			if target.Arch != tt.wantArch {
				t.Errorf("Arch = %q, want %q", target.Arch, tt.wantArch)
			}
			if target.OS != tt.wantOS {
				t.Errorf("OS = %q, want %q", target.OS, tt.wantOS)
			}
			if target.Env != tt.wantEnv {
				t.Errorf("Env = %q, want %q", target.Env, tt.wantEnv)
			}
		})
	}
}

func TestTargetTripleString(t *testing.T) {
	tests := []struct {
		name   string
		target *TargetTriple
		want   string
	}{
		{
			name:   "with vendor",
			target: &TargetTriple{Arch: "x86_64", Vendor: "pc", OS: "windows", Env: "gnu"},
			want:   "x86_64-pc-windows-gnu",
		},
		{
			name:   "without vendor",
			target: &TargetTriple{Arch: "x86_64", Vendor: "unknown", OS: "linux", Env: "gnu"},
			want:   "x86_64-linux-gnu",
		},
		{
			name:   "without env",
			target: &TargetTriple{Arch: "aarch64", Vendor: "apple", OS: "darwin", Env: "none"},
			want:   "aarch64-apple-darwin-none",
		},
		{
			name:   "minimal",
			target: &TargetTriple{Arch: "x86_64", Vendor: "unknown", OS: "linux", Env: "gnu"},
			want:   "x86_64-linux-gnu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseTargetRoundTrip(t *testing.T) {
	triples := []string{
		"linux-x64",
		"linux-arm64",
		"windows-x64",
		"macos-arm64",
		"x86_64-unknown-linux-gnu",
		"aarch64-unknown-linux-musl",
	}

	for _, triple := range triples {
		t.Run(triple, func(t *testing.T) {
			target, err := ParseTarget(triple)
			if err != nil {
				t.Fatalf("ParseTarget(%q) failed: %v", triple, err)
			}

			result := target.String()
			target2, err := ParseTarget(result)
			if err != nil {
				t.Fatalf("ParseTarget(%q) failed: %v", result, err)
			}

			if target.Arch != target2.Arch || target.OS != target2.OS || target.Env != target2.Env {
				t.Errorf("Round trip failed: %q -> %q -> %q", triple, result, target2.String())
			}
		})
	}
}

func TestDetectHost(t *testing.T) {
	host := DetectHost()

	if host == nil {
		t.Fatal("DetectHost() returned nil")
	}

	if host.Arch == "" {
		t.Error("DetectHost() returned empty Arch")
	}

	if host.OS == "" {
		t.Error("DetectHost() returned empty OS")
	}

	t.Logf("Host triple: %s", host.String())
}

func TestCommonTargets(t *testing.T) {
	for name, target := range CommonTargets {
		t.Run(name, func(t *testing.T) {
			if target.Arch == "" {
				t.Errorf("CommonTargets[%q] has empty Arch", name)
			}
			if target.OS == "" {
				t.Errorf("CommonTargets[%q] has empty OS", name)
			}
		})
	}
}
