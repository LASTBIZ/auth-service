package password

type Service struct {
	storage Storage
}

func NewPasswordService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}
