package dashboard

import (
	"html/template"
	"strings"
	"testing"
	"time"
)

func TestFormatTimeAgo(t *testing.T) {
	tests := []struct {
		name     string
		offset   time.Duration
		expected string
	}{
		{
			name:     "just now (0 seconds)",
			offset:   0,
			expected: "just now",
		},
		{
			name:     "just now (30 seconds)",
			offset:   30 * time.Second,
			expected: "just now",
		},
		{
			name:     "just now (59 seconds)",
			offset:   59 * time.Second,
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			offset:   1*time.Minute + time.Second,
			expected: "1m ago",
		},
		{
			name:     "5 minutes ago",
			offset:   5*time.Minute + time.Second,
			expected: "5m ago",
		},
		{
			name:     "59 minutes ago",
			offset:   59*time.Minute + time.Second,
			expected: "59m ago",
		},
		{
			name:     "1 hour ago",
			offset:   1*time.Hour + time.Second,
			expected: "1h ago",
		},
		{
			name:     "23 hours ago",
			offset:   23*time.Hour + time.Second,
			expected: "23h ago",
		},
		{
			name:     "1 day ago",
			offset:   24*time.Hour + time.Second,
			expected: "1d ago",
		},
		{
			name:     "2 days ago",
			offset:   48*time.Hour + time.Second,
			expected: "2d ago",
		},
		{
			name:     "30 days ago",
			offset:   30*24*time.Hour + time.Second,
			expected: "30d ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := time.Now().Add(-tt.offset)
			got := FormatTimeAgo(ts)
			if got != tt.expected {
				t.Errorf("FormatTimeAgo(-%v) = %q, want %q", tt.offset, got, tt.expected)
			}
		})
	}
}

func TestMapActivityItem(t *testing.T) {
	renderIcon := func(name string) template.HTML {
		return template.HTML("<svg>" + name + "</svg>")
	}

	now := time.Now()

	tests := []struct {
		name             string
		eventType        string
		eventName        string
		wantTitle        string
		wantDescription  string
		wantIconContains string
	}{
		{
			name:             "user_created event",
			eventType:        "user_created",
			eventName:        "John Doe",
			wantTitle:        "New User Created",
			wantDescription:  "John Doe added",
			wantIconContains: "icon-user-plus",
		},
		{
			name:             "role_modified event",
			eventType:        "role_modified",
			eventName:        "Admin",
			wantTitle:        "Role Updated",
			wantDescription:  "Admin role modified",
			wantIconContains: "icon-shield",
		},
		{
			name:             "unknown event type falls back to generic",
			eventType:        "something_else",
			eventName:        "Test Entity",
			wantTitle:        "Activity",
			wantDescription:  "Test Entity",
			wantIconContains: "icon-info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := MapActivityItem(tt.eventType, tt.eventName, now, renderIcon)

			if item.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", item.Title, tt.wantTitle)
			}
			if item.Description != tt.wantDescription {
				t.Errorf("Description = %q, want %q", item.Description, tt.wantDescription)
			}
			if !strings.Contains(string(item.IconHTML), tt.wantIconContains) {
				t.Errorf("IconHTML = %q, want to contain %q", item.IconHTML, tt.wantIconContains)
			}
			if item.TimeAgo == "" {
				t.Error("TimeAgo should not be empty")
			}
		})
	}
}

// TestBuildGetUsersByRoleID_NilDB and TestBuildGetDashboardData_NilDB removed:
// the raw-SQL builder functions have been migrated to espyna use cases at
// packages/espyna-golang/internal/application/usecases/service/dashboard/home/.
