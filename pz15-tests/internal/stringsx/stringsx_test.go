package stringsx

import "testing"

func TestClip(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		got := Clip("", 5)
		want := ""
		if got != want {
			t.Fatalf("Clip(\"\", 5) = %q; want %q", got, want)
		}
	})

	t.Run("max equals zero", func(t *testing.T) {
		got := Clip("hello", 0)
		want := ""
		if got != want {
			t.Fatalf("Clip(\"hello\", 0) = %q; want %q", got, want)
		}
	})

	t.Run("max less than zero", func(t *testing.T) {
		got := Clip("hello", -3)
		want := ""
		if got != want {
			t.Fatalf("Clip(\"hello\", -3) = %q; want %q", got, want)
		}
	})

	t.Run("max equals length", func(t *testing.T) {
		s := "hello"
		got := Clip(s, len(s))
		want := s
		if got != want {
			t.Fatalf("Clip(%q, %d) = %q; want %q", s, len(s), got, want)
		}
	})

	t.Run("max greater than length", func(t *testing.T) {
		s := "hello"
		got := Clip(s, len(s)+10)
		want := s
		if got != want {
			t.Fatalf("Clip(%q, %d) = %q; want %q", s, len(s)+10, got, want)
		}
	})
}

func TestClip_UTF8(t *testing.T) {
	t.Run("unicode string", func(t *testing.T) {
		s := "Привет"
		got := Clip(s, 4)
		want := s[:4]
		if got != want {
			t.Fatalf("Clip(%q, 4) = %q; want %q", s, got, want)
		}
	})
}
