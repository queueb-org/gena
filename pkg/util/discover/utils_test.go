package discover

import "testing"

func TestMakePath(t *testing.T) {
	for _, test := range []struct {
		name     string
		fullPath string
		appName  string
		expected string
	}{
		{
			name:     "ok",
			fullPath: "/home/user/go/src/github.com/user/app/pkg/config/v1",
			appName:  "app",
			expected: "/home/user/go/src/github.com/user/app",
		},
		{
			name:     "not-found",
			fullPath: "/home/user/go/src/github.com/user/app/pkg/config/v1",
			appName:  "example-fqdn",
			expected: "/home/user/go/src/github.com/user/app/pkg/config/v1",
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := makePath(test.fullPath, test.appName); result != test.expected {
				in.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
