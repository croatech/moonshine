package domain

import "testing"

func TestUser_ReachedNewLevel(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name: "level 1 with enough exp",
			user: &User{
				Level: 1,
				Exp:   100,
			},
			expected: true,
		},
		{
			name: "level 1 with exact exp",
			user: &User{
				Level: 1,
				Exp:   100,
			},
			expected: true,
		},
		{
			name: "level 1 with not enough exp",
			user: &User{
				Level: 1,
				Exp:   99,
			},
			expected: false,
		},
		{
			name: "level 2 with enough exp",
			user: &User{
				Level: 2,
				Exp:   200,
			},
			expected: true,
		},
		{
			name: "level 10 with enough exp",
			user: &User{
				Level: 10,
				Exp:   20000,
			},
			expected: true,
		},
		{
			name: "level 10 with not enough exp",
			user: &User{
				Level: 10,
				Exp:   19999,
			},
			expected: false,
		},
		{
			name: "level 11 (not in matrix)",
			user: &User{
				Level: 11,
				Exp:   100000,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.ReachedNewLevel()
			if result != tt.expected {
				t.Errorf("ReachedNewLevel() = %v, want %v", result, tt.expected)
			}
		})
	}
}








