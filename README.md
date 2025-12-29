# ViitorBot

Un simple bot de Discord desarrollado en Go que obtiene artículos aleatorios de Wikipedia cuando se ejecuta el comando "viitorbot hazlotuyo".

## Características

- Escucha el comando "viitorbot hazlotuyo" en canales de Discord
- Obtiene un artículo aleatorio de Wikipedia en español
- Responde con el título, un extracto del artículo y un enlace para más información

## Requisitos

- Go 1.24 o superior
- Un token de bot de Discord

## Instalación

1. Clona este repositorio:

   ```
   git clone https://github.com/tu-usuario/viitorbot.git
   cd viitorbot
   ```

2. Instala las dependencias:

   ```
   go mod tidy
   ```

3. Crea un archivo `.env` en la raíz del proyecto con el siguiente contenido:
   ```
   DISCORD_TOKEN=tu_token_aqui
   ```

## Ejecución

```
go run main.go
```

## Cómo obtener un token de Discord

1. Ve al [Portal de Desarrolladores de Discord](https://discord.com/developers/applications)
2. Crea una nueva aplicación
3. Ve a la sección "Bot" y haz clic en "Add Bot"
4. Copia el token
5. En la sección "Bot", asegúrate de habilitar "MESSAGE CONTENT INTENT" en "Privileged Gateway Intents"

## Uso

Una vez que el bot esté en ejecución y añadido a tu servidor, simplemente escribe:

```
viitorbot hazlotuyo
```

El bot responderá con un artículo aleatorio de Wikipedia.

## Búsqueda por fecha

Además de devolver artículos aleatorios, el bot admite búsquedas relacionadas con una fecha. Si el mensaje contiene una referencia temporal, el bot intentará obtener un artículo relacionado con ese día (usando el feed "On this day" de Wikipedia) y añadirá evidencia sobre por qué el artículo está relacionado con la fecha.

Soportes de fecha reconocidos:

- Palabras clave: `hoy`, `ayer`, `mañana` (también `manana`).
- Formatos de fecha: `YYYY-MM-DD`, `DD/MM/YYYY`, `DD-MM-YYYY`.

Ejemplos de uso:

```
viitorbot hazlotuyo hoy
viitorbot hazlotuyo 2025-12-30
viitorbot hazlotuyo 30/12/2025
```

Qué devuelve el bot cuando se detecta una fecha:

- Título del artículo
- Extracto corto
- Enlace al artículo en Wikipedia
- Evidencia (provenance) extra: tipo (evento/nacimiento/muerte/fiesta), año y texto descriptivo del ítem del feed

La evidencia te ayuda a verificar automáticamente por qué el artículo está asociado a esa fecha (por ejemplo, porque en esa fecha ocurrió un evento histórico, nació alguien, etc.).

## Licencia

[MIT](LICENSE)
