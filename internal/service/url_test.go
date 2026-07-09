package service

import (
	"errors"
	"testing"
	"time"

	"github.com/baydogan/lnk/internal/errs"
)

func TestParseExpiryEmpty(t *testing.T) {
	for _, in := range []string{"", "   ", "\t"} {
		got, err := parseExpiry(in)
		if err != nil {
			t.Fatalf("parseExpiry(%q) unexpected error: %v", in, err)
		}
		if got != nil {
			t.Fatalf("parseExpiry(%q) = %v, want nil", in, got)
		}
	}
}

func TestParseExpiryValid(t *testing.T) {
	cases := map[string]time.Duration{
		"30m":   30 * time.Minute,
		"1h":    time.Hour,
		"2h30m": 2*time.Hour + 30*time.Minute,
		"1.5h":  90 * time.Minute,
		"7d":    7 * 24 * time.Hour,
		"2w":    14 * 24 * time.Hour,
		"1d":    24 * time.Hour,
	}
	for in, want := range cases {
		before := time.Now().Add(want)
		got, err := parseExpiry(in)
		after := time.Now().Add(want)
		if err != nil {
			t.Fatalf("parseExpiry(%q) unexpected error: %v", in, err)
		}
		if got == nil {
			t.Fatalf("parseExpiry(%q) = nil, want non-nil", in)
		}
		if got.Before(before.Add(-time.Second)) || got.After(after.Add(time.Second)) {
			t.Fatalf("parseExpiry(%q) = %v, want ~%v", in, got, before)
		}
	}
}

func TestParseExpiryInvalid(t *testing.T) {
	for _, in := range []string{"0d", "-1h", "abc", "10", "d", "w", "1x", "0", "-3d", "1.5"} {
		got, err := parseExpiry(in)
		if !errors.Is(err, errs.ErrExpireFormat) {
			t.Fatalf("parseExpiry(%q) err = %v, want ErrExpireFormat", in, err)
		}
		if got != nil {
			t.Fatalf("parseExpiry(%q) = %v, want nil", in, got)
		}
	}
}
