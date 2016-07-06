package geckoboard

import "testing"

func TestMergeInConfig(t *testing.T) {
	c1 := Config{Key: "foobar"}

	c2 := Config{}
	c2.mergeIn(c1)
	if c2.Key != "foobar" {
		t.Fatalf("Api Key is missing")
	}
}
