package domain

type Location struct {
	Model
	Name     string `db:"name"`
	Slug     string `db:"slug"`
	Cell     bool   `db:"cell"`
	Inactive bool   `db:"inactive"`
	Image    string `db:"image"`
	ImageBg  string `db:"image_bg"`
}

const (
	WaywardPinesSlug = "wayward_pines"
	MoonshineSlug    = "moonshine"
)
