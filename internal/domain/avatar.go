package domain

type Avatar struct {
	Model
	Private bool   `db:"private"`
	Image   string `db:"image"`
}
