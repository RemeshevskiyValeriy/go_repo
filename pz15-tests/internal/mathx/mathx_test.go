package mathx

import "testing"

func TestSum_Table(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{2, 3, 5}, {10, -5, 5}, {0, 0, 0},
	}
	for _, c := range cases {
		got := Sum(c.a, c.b)
		if got != c.want {
			t.Fatalf("Sum(%d,%d)=%d; want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestDivide(t *testing.T) {
	t.Run("normal division", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil {
			t.Fatalf("Divide(10,2) unexpected error: %v", err)
		}
		if got != 5 {
			t.Fatalf("Divide(10,2) = %d; want 5", got)
		}
	})

	t.Run("divide by zero", func(t *testing.T) {
		got, err := Divide(10, 0)
		if err == nil {
			t.Fatalf("Divide(10,0) expected error, got nil")
		}
		if got != 0 {
			t.Fatalf("Divide(10,0) = %d; want 0 on error", got)
		}
	})

	t.Run("negative divisor", func(t *testing.T) {
		got, err := Divide(10, -2)
		if err != nil {
			t.Fatalf("Divide(10,-2) unexpected error: %v", err)
		}
		if got != -5 {
			t.Fatalf("Divide(10,-2) = %d; want -5", got)
		}
	})

	t.Run("negative dividend", func(t *testing.T) {
		got, err := Divide(-10, 2)
		if err != nil {
			t.Fatalf("Divide(-10,2) unexpected error: %v", err)
		}
		if got != -5 {
			t.Fatalf("Divide(-10,2) = %d; want -5", got)
		}
	})
}
