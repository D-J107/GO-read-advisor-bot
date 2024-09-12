package files

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"read-adviser-bot/lib/er"
	"read-adviser-bot/storage"
)

type Storage struct {
	basePath string // в какой папке храним файлы
}

const defaultPerm = 0774

func New(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = er.Wrap("Cant save page to File Storage", err)
	}()
	// складываем все файлы конкретного юзера в папку с его именем
	filePath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)
	createdFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = createdFile.Close()
	}()

	// преобразовывает page в формат gob и записывает в указанный файл
	if err := gob.NewEncoder(createdFile).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = er.Wrap("Cant pick random page from File!", err) }()
	pathToFolder := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(pathToFolder)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrorNoSavedPages
	}

	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(pathToFolder, file.Name()))
}

func (s *Storage) Remove(page *storage.Page) (err error) {
	fName, err := fileName(page)
	if err != nil {
		return er.Wrap("Cant remove file "+fName, err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)
	err = os.Remove(path)
	if err != nil {
		return er.Wrap("Cant remove file "+fName, err)
	}

	return nil
}

func (s *Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, er.Wrap("Cant check if file "+fName+" exists!", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, er.Wrap("Cant check if file "+fName+" exists!", err)
	}
	return true, nil
}

func (s *Storage) decodePage(filePath string) (page *storage.Page, err error) {
	defer func() { err = er.Wrap("Cant decode page!", err) }()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	baseName, err := p.Hash()
	return baseName, err
}
