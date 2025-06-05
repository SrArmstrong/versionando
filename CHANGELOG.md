# CHANGELOG

Toda modificacion important en este trabajo se documenta aquí.

# [1.2.0] - 2025-06-05
# Added
- Endpoint `/api/tasks` para crear tareas
- Middleware de autenticación con JWT

# Changed
- El modelo `User` ahora incluye `FechaNacimiento` y `PreguntaSecreta`

# [1.1.0] - 2025-05-22
# Added
- Registro de usuarios con Firestore
- Hashing de contraseñas con bcrypt

# [1.0.0] - 2025-05-10
# Added
- Estructura base del proyecto con Fiber
