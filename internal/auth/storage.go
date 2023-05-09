package auth

type Storage interface {
	CreatePassword(hash PasswordHash) error
	DeletePassword(userID uint32) error
	UpdatePassword(hash PasswordHash) error
	GetHash(userID uint32) (PasswordHash, error)
	CreateProvider(provider OAuthProvider) error
	DeleteProvider(userID uint32, provider string) error
	UpdateProvider(provider OAuthProvider) error
	GetProvider(userID uint32, provider string) (OAuthProvider, error)
}
