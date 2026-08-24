# Backend

API REST del sistema de gestión de biblioteca (LMS), encargada de la lógica de negocio: inventario de libros, registros de estudiantes, préstamos, devoluciones y penalizaciones.

## Tech Stack

* **Lenguaje:** Go (Golang)
* **Base de datos:** PostgreSQL

## Contenido

* `Dockerfile`: define la imagen del backend (build multi-stage con Go) usada por el `docker-compose.yml` de la raíz del proyecto.
* Código fuente de la API (a agregar).

## Desarrollo local

El backend se construye y levanta junto con el resto de servicios (frontend y base de datos) desde la raíz del repositorio:

```bash
docker compose up --build backend
```
