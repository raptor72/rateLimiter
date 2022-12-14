package whitelists

type Storage interface {
	Select() ([]*WhiteListModel, error)
	// SelectByAddress(address string) ([]*WhiteListModel, error)
	// Create(model *WhiteListModel) (*WhiteListModel, error)
	// Delete(id int) error
	// Update(model *UpdateWhiteListModel) (*WhiteListModel, error)
}

type WhiteListModel struct {
	ID      int    `db:"id" json:"id"`
	Address string `db:"address" json:"address"`

	Total uint `db:"total" json:"-"`
}

type UpdateWhiteListModel struct {
	ID      int     `db:"id" json:"id"`
	Address *string `db:"address" json:"address"`

	Total uint `db:"total" json:"-"`
}
