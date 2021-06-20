package db

const (
	sportList = "list"
)

func getSportQuery() map[string]string {
	return map[string]string{
		sportList: `
			SELECT 
				id, 
				sport_type, 
				name, 
				details,
				visible, 
				advertised_start_time
			FROM sports
		`,
	}
}
