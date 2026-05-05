package helm

import (
	"testing"
)

func TestReleaseRef_String_WithNamespace(t *testing.T) {
	ref := ReleaseRef{Name: "myapp", Namespace: "production"}
	got := ref.String()
	want := "production/myapp"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestReleaseRef_String_NoNamespace(t *testing.T) {
	ref := ReleaseRef{Name: "myapp", Namespace: ""}
	got := ref.String()
	want := "myapp"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestReleaseRef_Key_MatchesString(t *testing.T) {
	ref := ReleaseRef{Name: "svc", Namespace: "staging"}
	if ref.Key() != ref.String() {
		t.Errorf("Key() = %q, want same as String() = %q", ref.Key(), ref.String())
	}
}
