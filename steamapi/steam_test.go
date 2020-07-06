package steamapi

import (
	"testing"
)

func TestReadBinary(t *testing.T) {

	c := NewClient()

	_, err := c.GetAliases(76561197968626192)
	if err != nil {
		t.Error(err)
	}

	group, err := c.GetGroup("103582791434672565", "", 1)
	if err != nil {
		t.Error(err)
	}
	if group.Details.Name != "Steam Universe" {
		t.Error("name")
	}
}
