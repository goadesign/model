package expr

import "testing"

func TestPersonEvalName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, want string
	}{
		{name: "", want: "unnamed person"},
		{name: "foo", want: `person "foo"`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			person := Person{
				Element: &Element{
					Name: tt.name,
				},
			}
			if got := person.EvalName(); got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
