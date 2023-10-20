package jsonEdit

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	uuid "github.com/satori/go.uuid"
)

type Item struct {
	Id       uuid.UUID
	Path     string
	Folder   string
	Date     string
	Realized bool
}

type Items []Item

type JSONFile struct {
	FileName string
}

func (j *JSONFile) Add(data Item) (Items, error) {
	jsonRead, errRead := j.Read()
	if errRead != nil {
		jsonRead = make(Items, 0)
	}
	if data.Id == uuid.Nil {
		data.Id = getUniqueID()
		jsonRead = append(jsonRead, data)
	} else {
		return make(Items, 0), errors.New("Este item TIENE Id, use Update")
	}

	errWrite := j.writeJSON(jsonRead)
	if errWrite != nil {
		return make(Items, 0), errWrite
	}
	return jsonRead, nil
}

func (j *JSONFile) Update(Id string, data Item) (Items, error) {
	jsonRead, errRead := j.Read()
	if errRead != nil {
		jsonRead = make(Items, 0)
	}
	id, _ := uuid.FromString(Id)

	var itemFound *Item
	for i, it := range jsonRead {
		if it.Id == id {
			itemFound = &jsonRead[i]
			break
		}
	}
	*itemFound = data
	itemFound.Id = id

	errWrite := j.writeJSON(jsonRead)
	if errWrite != nil {
		return make(Items, 0), errWrite
	}
	return jsonRead, nil
}

func (j *JSONFile) Read() (Items, error) {
	file, err := os.Open(j.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, errReadAll := io.ReadAll(file)
	if errReadAll != nil {
		return nil, errReadAll
	}

	var data Items
	json.Unmarshal(byteValue, &data)

	return data, nil
}

func (j *JSONFile) Delete(Id string) (Items, error) {
	jsonRead, errRead := j.Read()
	if errRead != nil {
		return make(Items, 0), errors.New("No existe archivo")
	}

	id, _ := uuid.FromString(Id)

	index := -1
	for i, item := range jsonRead {
		if item.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		return make(Items, 0), errors.New("Item not found")
	}
	jsonRead = append(jsonRead[:index], jsonRead[index+1:]...)

	errWrite := j.writeJSON(jsonRead)
	if errWrite != nil {
		return make(Items, 0), errWrite
	}
	return jsonRead, nil
}

func getUniqueID() uuid.UUID {
	return uuid.NewV4()
}

func (j *JSONFile) writeJSON(data Items) error {
	jsonData, errConv := json.Marshal(data)
	if errConv != nil {
		return errConv
	}

	file, err := os.Create(j.FileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	return err
}
