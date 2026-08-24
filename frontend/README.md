# Frontend

Interfaz web del sistema de gestión de biblioteca (LMS), consumida por el administrador único del sistema para gestionar libros, estudiantes, préstamos, devoluciones y penalizaciones.

## Tech Stack

* **Framework:** React.js

## Contenido

* `Dockerfile`: define la imagen del frontend (build multi-stage con Node.js, servida con Nginx) usada por el `docker-compose.yml` de la raíz del proyecto.
* Código fuente de la aplicación (a agregar).

## Desarrollo local

El frontend se construye y levanta junto con el resto de servicios (backend y base de datos) desde la raíz del repositorio:

```bash
docker compose up --build frontend
```
