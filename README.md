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

## Peticiones HTTPS

### PUT

Endpoint esperado (Se remplaza el "6ywBh2GKBaq2crHAnTv1" por el usuario requerido): http://127.0.0.1:3000/api/users/6ywBh2GKBaq2crHAnTv1

**Campos permitidos para ser editados**
{
  "nombre": "Sergio",
  "apellidos": "Pérez Aldavalde",
  "email": "sergio@gmail.com",
  "password": "Ares"
}

### DELETE

Endpoint esperado (Se remplaza el "6ywBh2GKBaq2crHAnTv1" por el usuario requerido): http://127.0.0.1:3000/api/delete/6ywBh2GKBaq2crHAnTv1


## Ingresar una tarea JSON

{
  "titulo": "Estudiar para el examen",
  "descripcion": "Repasar los temas de la unidad 3 y 4",
  "fecha_limite": "2025-06-10T23:59:59Z"
}
