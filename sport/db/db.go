package db

import (
	"fmt"
	"log"
	"time"

	"syreclabs.com/go/faker"
)


func (r *sportsRepo) seed(number int) error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS sports 
		(id INTEGER PRIMARY KEY, 
		sport_type TEXT,
		meeting_id INTEGER, 
		name TEXT, 
		details TEXT,
		visible INTEGER, 
		advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	if (number <= 0) {
		log.Fatal("Number passed must be grater than 0");
		return nil
	}

	for i := 1; i <= number; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO sports(id, sport_type, name, details,visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			var typeStr = faker.RandomChoice(([]string{"boxing","cricket","racing","lawn bowls"}))
			_, err = statement.Exec(
				i,// id
				typeStr, // type
				fmt.Sprintf("%v No. %v",typeStr,i),// Name
				fmt.Sprintf("It's game %v of the %v event",i,typeStr), // Details
				faker.Number().Between(0, 1), // visible
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339), // Advertised Start time
			)
		}
	}

	return err
}

func (r *sportsRepo) defaultSeed() error {
	return r.seed(100)
}
