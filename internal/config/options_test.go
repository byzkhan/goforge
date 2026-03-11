package config

import "testing"

func TestServiceOptions_Validate(t *testing.T) {
	base := func() ServiceOptions {
		o := Defaults()
		o.Name = "myapi"
		o.Module = "github.com/org/myapi"
		o.OutputDir = "/tmp/myapi"
		return o
	}

	t.Run("valid defaults", func(t *testing.T) {
		if err := base().Validate(); err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		o := base()
		o.Name = ""
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for empty name")
		}
	})

	t.Run("invalid name", func(t *testing.T) {
		for _, bad := range []string{"My-Api", "123start", "has spaces", "UPPER"} {
			o := base()
			o.Name = bad
			if err := o.Validate(); err == nil {
				t.Fatalf("expected error for name %q", bad)
			}
		}
	})

	t.Run("empty module", func(t *testing.T) {
		o := base()
		o.Module = ""
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for empty module")
		}
	})

	t.Run("invalid router", func(t *testing.T) {
		o := base()
		o.Router = "gin"
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for unsupported router")
		}
	})
}
