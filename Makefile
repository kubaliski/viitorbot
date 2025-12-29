VERSION ?= v2.0.0

# Detect OS
ifeq ($(OS),Windows_NT)
    RM = del /q
    RMDIR = rmdir /s /q
    MKDIR = if not exist bin mkdir bin
    GOLINT = $(shell if exist tools\golangci-lint.exe (echo tools\golangci-lint.exe) else (echo golangci-lint))
else
    RM = rm -f
    RMDIR = rm -rf
    MKDIR = mkdir -p bin
    GOLINT = $(shell if [ -f ./tools/golangci-lint ]; then echo ./tools/golangci-lint; else echo golangci-lint; fi)
endif

# Docker variables
COMPOSE = docker compose
SERVICE = viitorbot

.PHONY: all example test cover tag build clean lint run help
.PHONY: docker-build docker-up docker-down docker-restart docker-logs docker-logs-f
.PHONY: docker-rebuild docker-clean docker-status docker-shell docker-deploy

# === Go commands ===
all: test

test:
	go test ./... -coverprofile=coverage.out

cover: test
	go tool cover -func=coverage.out

tag:
	git tag -a $(VERSION) -m "release $(VERSION)"
	git push origin $(VERSION)

run:
	go run ./cmd/viitorbot

build:
	$(MKDIR)
	go build -o bin/viitorbot ./cmd/viitorbot

lint:
	$(GOLINT) run --timeout 5m ./...

clean:
	$(RM) coverage.out
	$(RMDIR) bin

# === Docker commands ===
docker-build: ## Construir imagen Docker
	$(COMPOSE) build

docker-up: ## Levantar contenedor en segundo plano
	$(COMPOSE) up -d

docker-down: ## Parar y eliminar contenedor
	$(COMPOSE) down

docker-restart: docker-down docker-up ## Reiniciar contenedor

docker-logs: ## Ver logs estáticos
	$(COMPOSE) logs $(SERVICE)

docker-logs-f: ## Ver logs en tiempo real
	$(COMPOSE) logs -f $(SERVICE)

docker-rebuild: docker-down ## Reconstruir desde cero
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d
	@echo "Contenedor reconstruido y levantado. Usa 'make docker-logs-f' para ver logs."

docker-clean: docker-down ## Limpiar contenedores e imágenes
	docker rmi viitorbot-$(SERVICE) 2>/dev/null || true

docker-status: ## Ver estado del contenedor
	@$(COMPOSE) ps
	@echo ""
	@docker stats $(SERVICE) --no-stream 2>/dev/null || echo "Contenedor no está corriendo"

docker-shell: ## Abrir shell en el contenedor
	docker exec -it $(SERVICE) sh

docker-deploy: docker-rebuild docker-logs-f ## Desplegar cambios y ver logs

help: ## Mostrar ayuda
	@echo "=== Comandos Go ==="
	@echo "  make test          - Ejecutar tests"
	@echo "  make cover         - Ver cobertura de tests"
	@echo "  make run           - Ejecutar bot localmente"
	@echo "  make build         - Compilar binario"
	@echo "  make lint          - Ejecutar linter"
	@echo "  make clean         - Limpiar archivos generados"
	@echo "  make tag           - Crear tag de versión"
	@echo ""
	@echo "=== Comandos Docker ==="
	@echo "  make docker-build      - Construir imagen"
	@echo "  make docker-up         - Levantar contenedor"
	@echo "  make docker-down       - Parar contenedor"
	@echo "  make docker-restart    - Reiniciar contenedor"
	@echo "  make docker-logs       - Ver logs"
	@echo "  make docker-logs-f     - Ver logs en tiempo real"
	@echo "  make docker-rebuild    - Reconstruir desde cero"
	@echo "  make docker-status     - Ver estado y recursos"
	@echo "  make docker-shell      - Abrir shell en contenedor"
	@echo "  make docker-deploy     - Desplegar cambios (rebuild + logs)"
	@echo "  make docker-clean      - Limpiar todo"
