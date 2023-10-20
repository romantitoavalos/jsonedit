package jsonedit

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"reflect"

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

	index := findItem(jsonRead, Id)
	if index == -1 {
		return make(Items, 0), errors.New("El Item no existe!!")
	}

	UpdateStruct(Id, &jsonRead[index], data)

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

	index := findItem(jsonRead, Id)
	if index == -1 {
		return make(Items, 0), errors.New("El Item no existe!!")
	}
	jsonRead = append(jsonRead[:index], jsonRead[index+1:]...)

	errWrite := j.writeJSON(jsonRead)
	if errWrite != nil {
		return make(Items, 0), errWrite
	}
	return jsonRead, nil
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

func getUniqueID() uuid.UUID {
	return uuid.NewV4()
}

func findItem(json Items, Id string) int {
	id, _ := uuid.FromString(Id)
	index := -1
	for i, item := range json {
		if item.Id == id {
			index = i
			break
		}
	}
	return index
}

func UpdateStruct(Id string, orig, new interface{}) {
	id, _ := uuid.FromString(Id)

	o := reflect.ValueOf(orig).Elem()
	n := reflect.ValueOf(new)

	t := o.Type()
	fields := make([]string, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		fields[i] = t.Field(i).Name
	}

	for _, name := range fields {
		f := n.FieldByName(name)

		if f.IsValid() {
			o.FieldByName(name).Set(f)
		}
	}

	o.FieldByName("Id").Set(reflect.ValueOf(id))
}
