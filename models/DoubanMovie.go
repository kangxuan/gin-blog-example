package models

type DoubanMovie struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Other    string `json:"other"`
	Desc     string `json:"desc"`
	Year     string `json:"year"`
	Area     string `json:"area"`
	Tag      string `json:"tag"`
	Star     string `json:"star"`
	Comment  string `json:"comment"`
	Quote    string `json:"quote"`
}

// AddDoubanMovie 添加电影
func AddDoubanMovie(movies []DoubanMovie) bool {
	db.CreateInBatches(movies, 100)
	return true
}
