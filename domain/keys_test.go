package domain

import (
	"strings"
	"testing"
)

func TestHashKeyDeterministic(t *testing.T) {
	a := HashKey("lnk_example")
	if a != HashKey("lnk_example") {
		t.Fatal("HashKey not deterministic")
	}
	if len(a) != 64 {
		t.Fatalf("HashKey len = %d, want 64", len(a))
	}
	if HashKey("lnk_example") == HashKey("lnk_other") {
		t.Fatal("HashKey collided on distinct inputs")
	}
}

func TestHashKeyKnownVector(t *testing.T) {
	const want = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got := HashKey("hello"); got != want {
		t.Fatalf("HashKey(hello) = %q, want %q", got, want)
	}
}

func TestGenerateAPIKeyFormat(t *testing.T) {
	pt, hash, prefix, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey: %v", err)
	}
	if !strings.HasPrefix(pt, KeyPrefix) {
		t.Fatalf("plaintext %q missing prefix", pt)
	}
	if strings.Contains(pt, "=") {
		t.Fatalf("plaintext %q contains padding", pt)
	}
	if len(prefix) != 12 || prefix != pt[:12] {
		t.Fatalf("prefix %q != first 12 of %q", prefix, pt)
	}
	if hash != HashKey(pt) {
		t.Fatal("hash != HashKey(plaintext)")
	}
}

func TestGenerateAPIKeyUnique(t *testing.T) {
	seen := make(map[string]struct{})
	for range 100 {
		pt, _, _, err := GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey: %v", err)
		}
		if _, dup := seen[pt]; dup {
			t.Fatalf("duplicate %q", pt)
		}
		seen[pt] = struct{}{}
	}
}
