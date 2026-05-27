package steamapi

import (
	"context"
	"testing"
)

func TestPlayers(t *testing.T) {

	client, server := createMockClient(mockHandler(200, `[{"newname": "test", "timechanged": "2023-01-01"}]`))
	defer server.Close()

	_, _, err := client.GetAliases(context.Background(), 76561197968626192)
	if err != nil {
		t.Error(err)
	}
}

func TestGroups(t *testing.T) {

	client, server := createMockClient(mockHandler(200, `<memberList><groupID64>103582791434672565</groupID64><groupDetails><groupName>Steam Universe</groupName></groupDetails></memberList>`))
	defer server.Close()

	group, _, err := client.GetGroup(context.Background(), "103582791434672565", "", 1)
	if err != nil {
		t.Error(err)
	}
	if group.Details.Name != "Steam Universe" {
		t.Error("name")
	}
}
