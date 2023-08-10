package config

import "errors"

var (
	ErrNoDbURI  = errors.New("no database URI")
	ErrNoDbName = errors.New("no database name")
)

type Storage struct {
	URI    string
	DbName string
}

func (s *Storage) load(envPrefix string) error {
	v := setupViper(envPrefix)

	s.URI = v.GetString("uri")
	if s.URI == "" {
		return ErrNoDbURI
	}

	s.DbName = v.GetString("db.name")
	if s.DbName == "" {
		return ErrNoDbName
	}

	return nil
}
