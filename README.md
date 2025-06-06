# CRUD y Manejo de Versiones con Git

Práctica creada con el propósito de familiarizarse con el lenguaje de programación Go, implementando un sistema CRUD que ya se había desarrollado previamente en prácticas de cuatrimestres anteriores. Este proyecto incluye autenticación mediante JWT y validación de rutas.

## Inicializador del modulo

go mod init versionando

## Actualizar dependencias

go mod tidy

## Ejecutar el programa

go run main.go

## Forma de incertar un nuevo usuario JSON

{
  "nombre": "Sergio",
  "apellidos": "Pérez Aldavalde",
  "email": "sergio@gmail.com",
  "password": "Ares",
  "fecha_nacimiento": "2004-01-11T00:00:00Z",
  "pregunta_secreta": "¿Cuál es tu color favorito?",
  "respuesta_secreta": "Azul"
}

## Ingresar mediante el login JSON

{
  "email": "sergio@gmail.com",
  "password": "Ares"
}

## Peticiones HTTPS para usuario "user"

### Get todos los usuario

Endpoint esperado: http://127.0.0.1:3000/api/users/

Salida:

[
  {
    "id": "h72pop7K32civl2dQtkO",
    "nombre": "Sergio",
    "apellidos": "Perez Aldavalde",
    "email": "sergio@gmail.com",
    "createdAt": "2025-06-05T18:24:08.72618Z",
    "updatedAt": "2025-06-05T18:24:08.72618Z"
  }
]

### Get usuario especifico

Endpoint esperado (Se remplaza el ":id" por el id del usuario requerido): http://127.0.0.1:3000/api/users/:id

Salida:

{
  "id": "h72pop7K32civl2dQtkO",
  "nombre": "Sergio",
  "apellidos": "Perez Aldavalde",
  "email": "sergio@gmail.com",
  "password": "",
  "fecha_nacimiento": "2004-01-11T00:00:00Z",
  "pregunta_secreta": "¿Cuál es tu color favorito?",
  "respuesta_secreta": "",
  "created_at": "2025-06-05T18:24:08.72618Z",
  "updated_at": "2025-06-05T18:24:08.72618Z"
}

### PUT

Endpoint esperado (Se remplaza el ":id" por el id del usuario requerido): http://127.0.0.1:3000/api/users/:id

**Campos permitidos para ser editados**
{
  "nombre": "",
  "apellidos": "",
  "email": "",
  "password": ""
}

Salida:

{
  "Usuario actualizado exitosamente"
}

### DELETE

Endpoint esperado (Se remplaza el ":id" por el id del usuario requerido): http://127.0.0.1:3000/api/delete/:id

Salida:

{
  "Usuario eliminado exitosamente"
}


## Ingresar una tarea JSON

{
  "titulo": "Estudiar para el examen",
  "descripcion": "Repasar los temas de la unidad 3 y 4",
  "fecha_limite": "2025-06-10T23:59:59Z"
}

## Peticiones HTTPS para tareas "task"


### Get para obtener todas las tareas del usuario

Endpoint esperado: http://127.0.0.1:3000/api/tasks/

Salida:

[
  {
    "id": "3kfld83jf93jfl",
    "titulo": "Estudiar para el examen",
    "descripcion": "Repasar los temas de la unidad 3 y 4",
    "fecha_limite": "2025-06-10T23:59:59Z",
    "created_at": "2025-06-05T19:15:00Z"
  }
]

### Get para obtener tarea especifica

Endpoint esperado (Se remplaza :id por el id de la tarea): http://127.0.0.1:3000/api/tasks/:id

Salida:

{
  "id": "3kfld83jf93jfl",
  "titulo": "Estudiar para el examen",
  "descripcion": "Repasar los temas de la unidad 3 y 4",
  "fecha_limite": "2025-06-10T23:59:59Z",
  "created_at": "2025-06-05T19:15:00Z"
}


### Put actualizar tarea

Endpoint esperado (Se remplaza :id por el id de la tarea): http://127.0.0.1:3000/api/tasks/:id

Entrada y campos actualizables:

{
  "titulo": "Estudiar Go",
  "descripcion": "Repasar estructuras de datos",
  "fecha_limite": "2025-06-11T23:59:59Z"
}


Salida:

{
  "message": "Tarea actualizada exitosamente"
}

### Delete eliminar tarea

Endpoint esperado (Se remplaza :id por el id de la tarea): http://127.0.0.1:3000/api/tasks/:id

Salida:

{
  "message": "Tarea eliminada exitosamente"
}
