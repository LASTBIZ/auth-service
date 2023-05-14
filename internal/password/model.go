package password

import (
	"github.com/devfeel/mapper"
	"lastbiz/auth-service/pkg/errors"
)

type Hash struct {
	ID     uint32 `mapper:"id"`
	UserID uint32 `mapper:"user_id"`
	Hash   string `mapper:"hash"`
}

func (p *Hash) ToMap() (map[string]interface{}, error) {
	hashMap := make(map[string]interface{})
	err := mapper.AutoMapper(p, &hashMap)
	if err != nil {
		return hashMap, errors.Wrap(err, "mapper.Decode(password_hash)")
	}

	return hashMap, nil
}
