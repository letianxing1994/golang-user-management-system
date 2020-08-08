package config

import "gopkg.in/ini.v1"

type IniParser struct {
	confReader *ini.File //config reader
}

type IniParserError struct {
	errorInfo string
}

func (e *IniParserError) Error() string {
	return e.errorInfo
}

func (p *IniParser) Load(configFileName string) error {
	conf, err := ini.Load(configFileName)
	if err != nil {
		p.confReader = nil
		return err
	}
	p.confReader = conf
	return nil
}

func (p *IniParser) GetString(section string, key string) string {
	if p.confReader == nil {
		return ""
	}

	s := p.confReader.Section(section)
	if s == nil {
		return ""
	}

	return s.Key(key).String()
}

func (p *IniParser) GetInt32(section string, key string) int32 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	value, _ := s.Key(key).Int()

	return int32(value)
}

func (p *IniParser) GetUint32(section string, key string) uint32 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	value, _ := s.Key(key).Uint()

	return uint32(value)
}

func (p *IniParser) GetInt64(section string, key string) int64 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	value, _ := s.Key(key).Int64()
	return value
}

func (p *IniParser) GetUint64(section string, key string) uint64 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	value, _ := s.Key(key).Uint64()
	return value
}

func (p *IniParser) GetFloat32(section string, key string) float32 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	value, _ := s.Key(key).Float64()
	return float32(value)
}

func (p *IniParser) GetFloat64(section string, key string) float64 {
	if p.confReader == nil {
		return 0
	}

	s := p.confReader.Section(section)
	if s == nil {
		return 0
	}

	value, _ := s.Key(key).Float64()
	return value
}
