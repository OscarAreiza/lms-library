# Infra

Esta carpeta contendrá la configuración de infraestructura del proyecto (fuera del build de las imágenes, que se gestiona con el `docker-compose.yml` en la raíz del repo).

Aquí se irán agregando, a medida que el proyecto avance, elementos como:

* Configuración de despliegue (VPS, cloud provider, etc.)
* Scripts de aprovisionamiento / Infraestructura como Código (Terraform, Ansible, etc.)
* Configuración de reverse proxy / balanceador (Nginx, Traefik, etc.)
* Pipelines de CI/CD
* Variables y secretos de entornos (staging, producción)

Por ahora esta carpeta está vacía; su contenido se definirá más adelante.
