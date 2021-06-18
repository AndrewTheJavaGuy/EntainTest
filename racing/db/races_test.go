package db_test

import (
	"database/sql"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var racesRepo db.RacesRepo

func TestMain(m *testing.M) {

	racingDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	racesRepo = db.NewRacesRepo(racingDB)
	if err := racesRepo.Init(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	racingDB.Close()
	os.Exit(code)
}


func TestEmptyFilter(t *testing.T) {
	request := new(racing.ListRacesRequestFilter)

	races,_ := racesRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}
}

func TestSingleIdFilter(t *testing.T) {
	request := new(racing.ListRacesRequestFilter)
	request.MeetingIds = append(request.MeetingIds, 1)

	races,_ := racesRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}
}

func TestActiveList(t *testing.T) {
	request := new(racing.ListRacesRequestFilter)
	request.OnlyVisible = true

	races,_ := racesRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}
}

func TestInActiveList(t *testing.T) {
	request := new(racing.ListRacesRequestFilter)
	request.OnlyVisible = false

	races,_ := racesRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}
}

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

func TestOnlyVisibleAppears(t *testing.T) {
	request := new(racing.ListRacesRequestFilter)
	request.OnlyVisible = true

	races,_ := racesRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}

	for _, race := range races {
		if !race.GetVisible() {
			t.Fatal("Found a non-visible status when this should be only visible")
		}
	}

}

