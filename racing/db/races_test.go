package db_test

import (
	"git.neds.sh/matty/entain/racing/db"
	"testing"
)

func TestListNoFieldOrOrder(t *testing.T) {
	var order []string
	order = db.Add_order(order,"","")

	if len(order) != 0 {
		t.Fatal("Did not expect this to add to the order list when a field and order is missing")
	}
}

func TestListNoField(t *testing.T) {
	var order []string
	order = db.Add_order(order,"","asc")

	if len(order) != 0 {
		t.Fatal("Did not expect this to add to the order list when a field and order is missing")
	}
}

func TestListNoOrder(t *testing.T) {
	var order []string
	order = db.Add_order(order,"test_field","")

	if len(order) != 0 {
		t.Fatal("Did not expect this to add to the order list when a field and order is missing")
	}
}

func TestListGoodOrder(t *testing.T) {
	var order []string
	order = db.Add_order(order,"test_field","asc")

	if len(order) != 1 {
		t.Fatal("Expected this to populate 1 row. Instead got ",len(order))
	}
}

func TestListBadOrder(t *testing.T) {
	var order []string
	order = db.Add_order(order,"test_field","Blah")

	if len(order) != 0 {
		t.Fatal("Expected this to populate 1 row. Instead got ",len(order))
	}
}

