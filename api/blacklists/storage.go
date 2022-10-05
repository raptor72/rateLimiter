package blacklists

type Storage interface {
	Select() ([]*BlackListModel, error)
	// SelectByAddress(address string) ([]*BlackListModel, error)
	// Create(model *BlackListModel) (*BlackListModel, error)
	// Delete(id int) error
	// Update(model *UpdateBlackListModel) (*BlackListModel, error)
}

type BlackListModel struct {
	ID      int    `db:"id" json:"id"`
	Address string `db:"address" json:"address"`

	Total uint `db:"total" json:"-"`
}

type UpdateBlackListModel struct {
	ID      int     `db:"id" json:"id"`
	Address *string `db:"address" json:"address"`

	Total uint `db:"total" json:"-"`
}
