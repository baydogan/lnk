package service

import (
	"strings"
	"testing"
)

func TestHashKeyDeterministic(t *testing.T) {
	a := hashKey("lnk_example")
	b := hashKey("lnk_example")
	if a != b {
		t.Fatalf("hashKey not deterministic: %q != %q", a, b)
	}
	if len(a) != 64 {
		t.Fatalf("hashKey len = %d, want 64", len(a))
	}
	if hashKey("lnk_example") == hashKey("lnk_other") {
		t.Fatal("hashKey collided on distinct inputs")
	}
}

func TestHashKeyKnownVector(t *testing.T) {
	const want = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got := hashKey("hello"); got != want {
		t.Fatalf("hashKey(hello) = %q, want %q", got, want)
	}
}

func TestGenerateAPIKeyFormat(t *testing.T) {
	pt, hash, prefix, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey error: %v", err)
	}
	if !strings.HasPrefix(pt, keyPrefix) {
		t.Fatalf("plaintext %q missing prefix %q", pt, keyPrefix)
	}
	if strings.Contains(pt, "=") {
		t.Fatalf("plaintext %q contains base64 padding", pt)
	}
	if len(prefix) != 12 || prefix != pt[:12] {
		t.Fatalf("prefix %q != first 12 of %q", prefix, pt)
	}
	if hash != hashKey(pt) {
		t.Fatalf("hash %q != hashKey(plaintext)", hash)
	}
}

func TestGenerateAPIKeyUnique(t *testing.T) {
	seen := make(map[string]struct{})
	for range 100 {
		pt, _, _, err := generateAPIKey()
		if err != nil {
			t.Fatalf("generateAPIKey error: %v", err)
		}
		if _, dup := seen[pt]; dup {
			t.Fatalf("generateAPIKey produced duplicate %q", pt)
		}
		seen[pt] = struct{}{}
	}
}
