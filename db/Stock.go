package db

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
	"structs"

	"go.etcd.io/bbolt"
)

var (
	ErrDataBase = errors.New("data base err")

	ErrNotFoundInStock            = errors.New("not found in Stock")
	ErrContainerIsExistInStock    = errors.New("is exist")
	ErrContainerIsNotExistInStock = errors.New("is not exist")
	ErrKindNotFound               = errors.New("cant found the kind")
	ErrKindIsExited               = errors.New(" kind is exited ")

	ErrModelSize         = errors.New("model size is empty or Colors is empty")
	ErrModelPrice        = errors.New("Price is Not acceptable")
	ErrModelDescription  = errors.New("Description is Not acceptable")
	ErrModelLinkesImage  = errors.New("linkes Image is empty")
	ErrModelSizeNotFound = errors.New("Size not found")

	ErrModelCommintEmpty = errors.New("Commint is not accptable")
	ErrModelCommintStars = errors.New("Stars is not accptable")
)

type stock struct {
	pathStock string
	database  map[string]map[string]*bbolt.DB
}



func (s *stock) addCommint(id int, Container, kind string, commint structs.UserCommint) error {
	if len(commint.Commint) == 0 {
		return ErrModelCommintEmpty
	}
	if commint.Stars <= 0 {
		return ErrModelCommintStars
	}

	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.Commint = append(model.Commint, commint)

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpateLinkesImage(id int, Container, kind string, linkesImage []string) error {
	if len(linkesImage) == 0 {
		return ErrModelLinkesImage
	}
	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.LinkesImage = linkesImage

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpateDescription(id int, Description, Container, kind string) error {
	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.Description = Description

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpatePrice(id, Price int, Container, kind string) error {
	if Price <= 0 {
		return ErrModelPrice
	}
	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.Price = Price

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) GetModleSize(id int, Container, kind, sizeName string) (*structs.Size, error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return nil, errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	var size structs.Size
	err := v.View(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}
		v, ok := model.Sizes[sizeName]
		if !ok {
			return ErrModelSizeNotFound
		}
		size = v
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &size, nil
}
func (s *stock) DeleteSizeFromModel(id int, Container, kind, sizeName string) error {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}
		delete(model.Sizes, sizeName)

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpdataSizeFromModel(id int, Container, kind, sizeName string, size structs.Size) error {

	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}
		if len(model.Sizes) == 0 || model.Sizes == nil {
			model.Sizes = make(map[string]structs.Size)
		}
		model.Sizes[sizeName] = size
		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) AddModelToKind(Container, kind string, model *structs.Model, OutStock bool) error {

	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	if OutStock {
		if len(model.Sizes) == 0 {
			return ErrModelSize
		}
		for _, v := range model.Sizes {
			if len(v.Colors) == 0 {
				return ErrModelSize
			} else {
				for _, s := range v.Colors {
					if len(s.ColorName) == 0 || s.Qty <= 0 {
						return ErrModelSize
					}
				}
			}
		}

		if model.Price <= 0 {
			return ErrModelPrice
		}
		if len(model.Description) == 0 {
			return ErrModelDescription
		}
		if len(model.LinkesImage) == 0 {
			return ErrModelLinkesImage
		}
	}
	return v.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}
		if OutStock {

			c, err := b.NextSequence()
			if err != nil {
				// this not to clinc
				return err
			}
			model.Id = int(c)
		}
		data, err := json.Marshal(&model)
		if err != nil {
			// this not to clinc
			return err
		}
		return b.Put(itob(model.Id), data)
	})

}
func (s *stock) NewKindtoContainer(Container, kind string) error {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	_, ok := s.database[Container][kind]
	if ok {
		return ErrKindIsExited
	}

	b, err := bbolt.Open(s.pathStock+"/"+Container+"/"+kind, 0600, nil)
	if err != nil {
		return err
	}

	if err = b.Batch(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte(kind))
		return err
	}); err != nil {
		return err
	}
	s.database[Container][kind] = b

	return nil
}

func (s *stock) DeleteModelFromKind(Container, kind string, id int) error {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	return v.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}
		return b.Delete(itob(id))

	})

}
func (s *stock) GetModelsInKind(formId, count int, Container, kind string) (models []structs.Model, _ error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return models, ErrContainerIsNotExistInStock
	}
	err := v.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		// check if bucket is exit
		if b == nil {
			return ErrKindNotFound
		}
		c := b.Cursor()
		k, nextVuale := c.Seek(itob(formId))
		if k == nil {
			return ErrKindNotFound
		}
		model := structs.Model{}
		err := json.Unmarshal(nextVuale, &model)
		if err != nil {
			return err
		}
		models = append(models, model)

		k, nextVuale = c.Next()
		i := 0

		for k != nil && i < count {
			model := structs.Model{}
			err := json.Unmarshal(nextVuale, &model)
			if err != nil {
				return err
			}

			models = append(models, model)

			k, nextVuale = c.Next()
			i++
		}

		return nil
	})
	if err != nil {
		return models, nil
	}

	return models, nil
}

func (S *stock) GetAllContainer() (Container []string) {
	for k := range S.database {
		Container = append(Container, k)
	}
	return Container
}

func (S *stock) GetAllContainerAndKind() (ContainerAndKind map[string][]string) {
	ContainerAndKind = make(map[string][]string)
	for k, v := range S.database {
		ContainerAndKind[k] = nil
		for k := range v {
			ContainerAndKind[k] = append(ContainerAndKind[k], k)
		}
	}
	return ContainerAndKind
}

func (db *stock) IsExist(name string) bool {
	name = strings.TrimSpace(name)
	_, ok := db.database[name]
	return ok
}
func (db *stock) DeleteContainer (name string ) error {
	 err:=  os.RemoveAll(db.pathStock+"/"+name)
	 if err != nil {
		  return err
	 }
	 return nil 
}

func (db *stock) AddNewContainer(name string) error {
	name = strings.TrimSpace(name)
	if db.IsExist(name) {
		return ErrContainerIsExistInStock
	}
	db.database[name] = nil
	err := os.Mkdir(db.pathStock+"/"+name, 0770)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (db *stock) init() error {

	f, err := os.Open(db.pathStock)
	defer f.Close()
	if err != nil {
		return err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, v := range files {
		if v.IsDir() {
			db.database[v.Name()] = make(map[string]*bbolt.DB)

			SubDir, err := os.Open(db.pathStock + "/" + v.Name())

			if err != nil {
				return err
			}
			defer SubDir.Close()
			filesInSubDir, _ := SubDir.Readdir(0)
			if filesInSubDir == nil {
				continue
			}
			for _, file := range filesInSubDir {

				if !file.IsDir() && path.Base(v.Name()) == ".db" {
					lenCut := len(file.Name()) - 3
					ptrdb, err := bbolt.Open(db.pathStock+"/"+v.Name()+"/"+file.Name(), 0600, nil)
					if err != nil {
						return err
					}
					db.database[v.Name()][file.Name()[:lenCut]] = ptrdb

				}

			}

		}
	}

	return nil
}

func (s *stock) ContainerAndKindIsExited(Container, kind string) (*bbolt.DB, error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return nil, errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	return v, nil
}
