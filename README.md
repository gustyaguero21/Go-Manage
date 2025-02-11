# Go-Manage
Go-Manage

# CRUD de Usuarios en Go

Este proyecto es una aplicación CRUD (Alta, Baja, Modificación y Consulta) desarrollada en Go para registrar usuarios. Se ha creado con el objetivo de ayudar a comprender el uso de **mocks**, especialmente con `sqlmock`.

## Características

- 🚀 Implementación de CRUD de usuarios.
- 🛠 Uso de SQLite3 como base de datos.
- ✅ Mejores prácticas de programación.
- 🔍 Pruebas unitarias con `sqlmock`.

## Tecnologías Utilizadas

- **Go** - Lenguaje principal.
- **SQLite3** - Base de datos ligera.
- **sqlmock** - Simulación de consultas SQL para pruebas.

## Instalación

Clonar el repositorio e instalar las dependencias:

```bash
# Clona el repositorio
git clone https://github.com/gustyaguero21/Go-Manage.git

# Entra en el directorio
cd tu-repo

# Instala sqlite3
go get github.com/mattn/go-sqlite3


# Instala las dependencias
go mod tidy
```



## Uso

Ejecutar la aplicación:

```bash
go run cmd/api/main.go
```

Ejecutar las pruebas con mocks:

```bash
go test ./...
```

## 📩 Colección de Postman

Puedes importar la colección de Postman desde el siguiente enlace:

[📥 Descargar colección de Postman](https://drive.google.com/file/d/1kL30kmvAYbBvWf1nRIgZBasr_dNlc77f/view?usp=sharing)

## Contribuir

Si deseas contribuir, por favor:
1. Realiza un fork del repositorio.
2. Crea una nueva rama (`feature/nueva-funcionalidad`).
3. Realiza los cambios y haz un commit (`git commit -m 'Agrega nueva funcionalidad'`).
4. Envía un pull request.

## Licencia

Este proyecto es de codigo abierto, cualquier persona que desee hacer uso de el, puede hacerlo. Tambien si quieren contribuir tanto en mejoras como correcciones, son mas que bienvenidos.

---

¡Gracias por visitar este proyecto! 🚀

