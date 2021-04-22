package helpers

import "testing"

func TestGetDocPath(t *testing.T) {
	t.Run("ensure correct doc path", func(t *testing.T) {
		got := GetDocPath("mycollection", "myitem")
		want := "mycollection/myitem"

		if got != want {
			t.Errorf("Wrong doc path. Got: %s; want: %s", got, want)
		}
	})
}
