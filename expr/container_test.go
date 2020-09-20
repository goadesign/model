package expr

import (
	"testing"
)

func TestContainerEvalName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{name: "", want: "unnamed container"},
		{name: "foo", want: `container "foo"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			container := Container{
				Element: &Element{
					Name: tt.name,
				},
			}
			if got := container.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
