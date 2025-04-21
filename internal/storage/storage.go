package storage

import "errors"

type Storage interface {
	InsertURL(uid string, url string) error
	GetURL(uid string) (string, error)
}

// тип urlStorage и его параметр Data
type UrlStorage struct {
	Data map[string]string
}

// конструктор объектов с типом urlStorage
func NewStorageStruct() *UrlStorage {
	return &UrlStorage{
		Data: make(map[string]string),
	}
}

// тип urlStorage и его метод InsertURL
func (s *UrlStorage) InsertURL(uid string, url string) error {
	s.Data[uid] = url
	return nil
}

// тип urlStorage и его метод GetURL
func (s *UrlStorage) GetURL(uid string) (string, error) {
	e, existss := s.Data[uid]
	if !existss {
		return uid, errors.New("URL with such id doesn`t exist")
	}
	return e, nil
}

//*******************************************************************
// Реализую интерфейс Storage

func MakeEntry(s Storage, uid string, url string) {
	s.InsertURL(uid, url)
}

func GetEntry(s Storage, uid string) (string, error) {
	e, err := s.GetURL(uid)
	return e, err
}

//********************************************************************
