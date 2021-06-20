package db

import (
	"database/sql"
	"fmt"
	"github.com/AndrewTheJavaGuy/entain/sport/proto/sport"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
	"time"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sport repository.
	Init() error
	// Init the repository, letting you set the number of records
	InitWithSeed(number int) error

	// List will return a list of races.
	List(filter *sport.ListSportsFilter) ([]*sport.SportDetails, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}


// NewRacesRepo creates a new sport repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sport repository dummy data, letting you specify the number of records
func (s *sportsRepo) InitWithSeed(number int) error {
	var err error

	s.init.Do(func() {
		if number > 0  {
			err = s.seed(number)
		} else {
			err = s.defaultSeed()
		}
	})

	return err
}

// Init prepares the sport repository dummy data.
func (s *sportsRepo) Init() error {
	return s.InitWithSeed(0)
}

// List all sports using the filter
func (r *sportsRepo) List(filter *sport.ListSportsFilter) ([]*sport.SportDetails, error) {

	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportQuery()[sportList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanSports(rows)
}

// Create the query from the filter
func (r *sportsRepo) applyFilter(query string, filter *sport.ListSportsFilter) (string, []interface{}) {
	var (
		clauses []string
		order_by   []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.SportType) > 0 {
		clauses = append(clauses, fmt.Sprintf("sport_type = '%v'",filter.SportType))
	}

	if filter.OnlyVisible {
		clauses = append(clauses, "visible = true")
	}

	if filter.Sort != nil {
		order_by = Add_order(order_by,"sport_type",filter.Sort.SportType)
		order_by = Add_order(order_by,"name",filter.Sort.Name)
		order_by = Add_order(order_by,"advertised_start_time",filter.Sort.AdvertisedStartTime)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	if (len(order_by) != 0) {
		query += " ORDER BY " + strings.Join(order_by, ", ")
	}

	return query, args
}

// If there's a asc or desc set, then add the value to the order by otherwise do nothing
func Add_order(order_by []string,field string,order string)([]string) {

	if len(field) == 0 || len(order) == 0 {
		return order_by
	}

	if order == "asc" || order == "desc" {
		order_by = append(order_by,field+" "+order)
	}

	return order_by
}

// Create the sportsdetails object to return from the query
func (m *sportsRepo) scanSports(
	rows *sql.Rows,
) ([]*sport.SportDetails, error) {
	var sports []*sport.SportDetails

	for rows.Next() {
		var sport sport.SportDetails
		var advertisedStart time.Time

		if err := rows.Scan(&sport.Id, &sport.SportType, &sport.Name, &sport.Details, &sport.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		if advertisedStart.Before(time.Now()) {
			sport.Status = "CLOSED"
		} else {
			sport.Status = "OPEN"
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		sport.AdvertisedStartTime = ts

		sports = append(sports, &sport)
	}

	return sports, nil
}
