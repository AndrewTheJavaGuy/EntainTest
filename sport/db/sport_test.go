package db_test

import (
	"database/sql"
	"github.com/AndrewTheJavaGuy/entain/sport/db"
	"github.com/AndrewTheJavaGuy/entain/sport/proto/sport"
	log "github.com/sirupsen/logrus"
	"os"
	"syreclabs.com/go/faker"
	"testing"
)

var totalRecords = 10

var sportsRepo db.SportsRepo

func TestMain(m *testing.M) {

	sportsDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	sportsRepo = db.NewSportsRepo(sportsDB)
	if err := sportsRepo.InitWithSeed(totalRecords); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	sportsDB.Close()
	os.Exit(code)
}


func TestEmptyFilter(t *testing.T) {
	request := new(sport.ListSportsFilter)
	sports,_ := sportsRepo.List(request)

	if len(sports) != totalRecords {
		t.Fatal("Expected to get",totalRecords,"back but instead got",len(sports))
	}
}

func TestActiveList(t *testing.T) {
	request := new(sport.ListSportsFilter)
	request.OnlyVisible = true

	races,_ := sportsRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}
}

func TestInActiveList(t *testing.T) {
	request := new(sport.ListSportsFilter)
	request.OnlyVisible = false

	races,_ := sportsRepo.List(request)

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
	request := new(sport.ListSportsFilter)
	request.OnlyVisible = true

	races,_ := sportsRepo.List(request)

	if len(races) == 0 {
		t.Fatal("Expected to get some results back")
	}

	for _, race := range races {
		if !race.GetVisible() {
			t.Fatal("Found a non-visible status when this should be only visible")
		}
	}
}

func TestOpenAndClosedStatus(t *testing.T) {
	request := new(sport.ListSportsFilter)

	races,returnError := sportsRepo.List(request)

	if returnError != nil {
		t.Fatal("Got an error back when not expected",returnError)
	}

	for _, race := range races {
		if race.GetStatus() != "OPEN" && race.GetStatus() != "CLOSED" {
			t.Fatal("Got an unexpected status back",race.GetStatus())
		}
	}

}

func TestSortOnAll(t *testing.T) {
	request := new(sport.ListSportsFilter)

	request.Sort = new(sport.Sort)
	request.Sort.SportType = faker.RandomChoice([]string{"asc","desc"})
	request.Sort.Name = faker.RandomChoice([]string{"asc","desc"})
	request.Sort.AdvertisedStartTime = faker.RandomChoice([]string{"asc","desc"})

	sports,returnError := sportsRepo.List(request)

	if returnError != nil {
		t.Fatal("Got an error back when not expected",returnError)
	}

	if len(sports) != totalRecords {
		t.Fatal("Expected to get",totalRecords,"back but instead got",len(sports))
	}
}
