package expr

import "testing"

func TestComponentEvalName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{name: "", want: "unnamed component"},
		{name: "foo", want: `component "foo"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			component := Component{
				Element: &Element{
					Name: tt.name,
				},
			}
			if got := component.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
