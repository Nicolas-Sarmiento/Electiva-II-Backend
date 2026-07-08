# Sistema de Monitoreo de pacientes en tiempo rea
Este repositorio contiene el código fuente para el backend de la asignatura Electiva II. Es una API desarrollada en **Go (Golang)** que integra un broker de mensajería **MQTT (Mosquitto)**, y está preparada para ejecutarse mediante contenedores para facilitar su despliegue y desarrollo local.


## Autores:
- Felipe Luna
- Nicolas Tinjaca
- Nicolas Sarmiento
  
## 🛠 Tecnologías y Herramientas

* **Lenguaje:** [Go (Golang)](https://go.dev/)
* **Mensajería / IoT:** [Eclipse Mosquitto (MQTT)](https://mosquitto.org/)
* **Contenedores:** [Docker](https://www.docker.com/) y [Docker Compose](https://docs.docker.com/compose/)
* **Arquitectura:** Estructura de directorios basada en los estándares de Go (`cmd/`, `internal/`)

## 📂 Estructura del Proyecto

El repositorio está organizado de la siguiente manera:

* `cmd/api/`: Contiene el punto de entrada principal de la aplicación (`main.go`). Es el encargado de inicializar y arrancar el servidor de la API.
* `internal/`: Contiene el núcleo y la lógica de negocio de la aplicación (modelos, controladores, repositorios, etc.) que no deben ser importados por otras aplicaciones.
* `api/`: Funciones, utilidades o especificaciones relacionadas directamente con el servicio web.
* `mosquitto/`: Configuraciones específicas para el broker MQTT (Mosquitto).
* `docker-compose.yml`: Archivo de orquestación de Docker que define y conecta los servicios necesarios (la API en Go, el broker MQTT, etc.).
* `endpoints.txt`: Documentación de referencia rápida con las rutas (endpoints) expuestas por la API.
* `go.mod` / `go.sum`: Archivos de configuración de Go para la gestión de dependencias del proyecto.

## ⚙️ Requisitos Previos

Para ejecutar este proyecto de forma local, es necesario contar con:

1. [Docker](https://docs.docker.com/get-docker/)
2. [Docker Compose](https://docs.docker.com/compose/install/)
3. *(Opcional)* [Go 1.x](https://go.dev/doc/install) en caso de que desees compilar el código o ejecutarlo fuera de Docker.

## 🚀 Instalación y Ejecución

La forma más ágil de levantar el entorno completo es utilizando Docker Compose. Esto inicializará la API backend y el broker Mosquitto en una red interna de Docker.

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/Nicolas-Sarmiento/Electiva-II-Backend.git](https://github.com/Nicolas-Sarmiento/Electiva-II-Backend.git)
   cd Electiva-II-Backend
   ```
2. Levantar los contenedores:
   ``` Bash
    docker-compose up --build
   ```

    Verificación de servicios:
      - La API debería arrancar e indicar por consola el puerto en el que está escuchando.
      -  El broker Mosquitto generalmente inicia en el puerto estándar MQTT 1883.

    Detener los servicios:
    Para detener y eliminar los contenedores generados, ejecuta:
    ```Bash
    docker-compose down
    ```
📡 Endpoints de la API

Las rutas de la API, los métodos HTTP requeridos y sus especificaciones se encuentran listados de manera detallada en el archivo endpoints.txt situado en la raíz del proyecto. Consulta ese documento para saber cómo realizar las peticiones al servidor.
📝 Notas de Desarrollo

   Configuración MQTT: Puedes ajustar el comportamiento del broker (puertos, autenticación, persistencia) modificando los archivos dentro de la carpeta /mosquitto.

   Recarga en caliente: Existe un directorio temporal (tmp_app) que normalmente sugiere el uso de herramientas como Air para live-reloading en Go. Si estás desarrollando activamente, los cambios podrían reflejarse en tiempo real.
