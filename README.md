## Este paquete proporciona funciones para leer, escribir, agregar, actualizar y eliminar elementos JSON desde un archivo:

### Instalación
```go
go get github.com/romantitoavalos/jsonedit
```

### Uso
```go
import "github.com/romantitoavalos/jsonedit"
```

### Crea una nueva instancia de JSONFile pasando el nombre del archivo:
```go 
json := jsonedit.JSONFile{FileName: "data.json"}
```

### Leer elementos
#### Usa el método Read() para obtener los elementos del archivo como un slice de structs Item:

```go
items, err := json.Read()
```

### Agregar elementos 

#### Usa el método Add() pasando un Item. Se le asignará un UUID automáticamente:
```go
item := jsonedit.Item{Path: "/home", Folder: "personal"}
items, err := json.Add(item) 
```

### Actualizar elementos

#### Usa el método Update() pasando el UUID del elemento a actualizar y los nuevos datos Item:
```go
updated := jsonedit.Item{Path: "/tmp"}
items, err := json.Update("179cc76a-5ba6-4b15-b3e0-b648d217c569", updated)
```

### Eliminar elementos
#### Usa el método Delete() pasando el UUID del elemento a eliminar:
```go
items, err := json.Delete("179cc76a-5ba6-4b15-b3e0-b648d217c569")
```

## Estructuras
#### El paquete define las siguientes estructuras:

- **Item**: representa un elemento con UUID, ruta, carpeta, fecha y flag realized
- **Items**: slice de Items
- **JSONFile**: contiene el nombre de archivo

## Errores
#### Los métodos pueden devolver los siguientes errores:

- Errores de lectura/escritura del archivo
- Item no encontrado al actualizar/eliminar
- UUID definido al intentar agregar un elemento nuevo

## Ejemplo
```go
json := jsonedit.JSONFile{FileName: "data.json"}

item1 := jsonedit.Item{Path: "/home/photos", Folder: "personal"}
items, _ := json.Add(item1)

item2 := jsonedit.Item{Path: "/etc", Folder: "system"} 
items, _ = json.Add(item2)

id := items[0].Id
updated := jsonedit.Item{Id: id, Realized: true}
items, _ = json.Update(id, updated) 

items, _ = json.Delete(id)

items, _ = json.Read()
```
