package steamapi

import (
	"testing"
)

func TestPlayers(t *testing.T) {

	c := NewClient()

	_, _, err := c.GetAliases(76561197968626192)
	if err != nil {
		t.Error(err)
	}
}

func TestGroups(t *testing.T) {

	c := NewClient()

	group, _, err := c.GetGroup("103582791434672565", "", 1)
	if err != nil {
		t.Error(err)
	}
	if group.Details.Name != "Steam Universe" {
		t.Error("name")
	}
}
