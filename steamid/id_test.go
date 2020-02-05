package steamid

import (
	"testing"
)

func TestGettersSetters(t *testing.T) {

	id := NewPlayerID(UniversePublic, AccountTypeIndividual, 1, 8360464)

	if id != 76561197968626192 {
		t.Error(id, 76561197968626192)
	}

	if id.GetUniverseID() != UniversePublic {
		t.Error("GetUniverseID", id.GetUniverseID(), UniversePublic)
	}

	if id.GetAccountType() != AccountTypeIndividual {
		t.Error("GetAccountType", id.GetAccountType(), AccountTypeIndividual)
	}

	if id.GetInstanceID() != 1 {
		t.Error("GetInstanceID", id.GetInstanceID(), 1)
	}

	if id.GetAccountID() != 8360464 {
		t.Error("GetAccountID", id.GetAccountID(), 4180232)
	}
}

func TestIDTypes(t *testing.T) {

	ids := []string{
		"STEAM_0:0:4180232",
		"U:1:8360464",
		"[U:1:8360464]",
		"U:1:8360464:1",
		"[U:1:8360464:1]",
		"8360464",
		"76561197968626192",
	}

	for _, idIn := range ids {

		id, err := ParsePlayerID(idIn)
		if err != nil {
			t.Error(err, idIn)
		}
		if id != 76561197968626192 {
			t.Error("not me", idIn)
		}
	}
}