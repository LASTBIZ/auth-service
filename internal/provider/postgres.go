package provider

import (
	sq "github.com/Masterminds/squirrel"
	storage2 "lastbiz/auth-service/internal/storage"
)

type Storage struct {
	queryBuilder sq.StatementBuilderType
	client       storage2.PostgreSQLClient
}

func NewProviderStorage(client storage2.PostgreSQLClient) *Storage {
	return &Storage{
		client:       client,
		queryBuilder: sq.StatementBuilder,
	}
}

const (
	scheme      = "public"
	table       = "passwords"
	tableScheme = scheme + "." + table
)

//func (s storage) CreateProvider(provider auth.OAuthProvider) error {
//	err := s.db.Create(&provider).Error
//	return err
//}
//
//func (s storage) UpdateProvider(provider auth.OAuthProvider) error {
//	err := s.db.
//		Save(provider).
//		Error
//	return err
//}
//
//func (s storage) DeleteProvider(userID uint32, provider string) error {
//	err := s.db.Model(&auth.OAuthProvider{}).
//		Where(&auth.OAuthProvider{UserID: userID, OAuthProvider: provider}).
//		Delete(&auth.OAuthProvider{}).
//		Error
//	return err
//}
//
//func (s storage) GetProvider(userID uint32, provider string) (auth.OAuthProvider, error) {
//	var _provider auth.OAuthProvider
//	err := s.db.Model(&auth.OAuthProvider{}).
//		Where(&auth.OAuthProvider{UserID: userID, OAuthProvider: provider}).
//		First(&_provider).
//		Error
//	return _provider, err
//}
