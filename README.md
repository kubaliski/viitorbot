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

## Licencia

[MIT](LICENSE)
