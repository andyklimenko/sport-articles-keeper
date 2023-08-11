package config

import "errors"

var ErrNoServerAddr = errors.New("no server address")

type Server struct {
	Address string
}

func (s *Server) load(envPrefix string) error {
	v := setupViper(envPrefix)

	s.Address = v.GetString("addr")
	if s.Address == "" {
		return ErrNoServerAddr
	}

	return nil
}
