package password

type Hash struct {
	ID     uint32 `mapper:"id"`
	UserID uint32 `mapper:"user_id"`
	Hash   string `mapper:"hash"`
}
