# Nome do executável que será gerado
APP_NAME = sudosee

# Caminho para o arquivo principal
MAIN_FILE = cmd/sudosee/main.go

# .PHONY avisa ao Make que essas palavras são comandos, e não nomes de pastas/arquivos
.PHONY: default run build clean fmt test

# Comando padrão ao digitar apenas 'make'
default: run

# Roda o projeto em modo de desenvolvimento
run:
	go run $(MAIN_FILE)

# Compila o projeto e gera um binário otimizado na pasta bin/
build:
	@echo "Compilando $(APP_NAME)..."
	go build -ldflags="-s -w" -o bin/$(APP_NAME) $(MAIN_FILE)
	@echo "Pronto! Executável gerado em bin/$(APP_NAME)"

# Remove a pasta de build
clean:
	@echo "Limpando arquivos compilados..."
	rm -rf bin/

# Formata todo o código fonte (Deixa tudo nos padrões oficiais do Go)
fmt:
	@echo "Formatando o código..."
	go fmt ./...

# Roda os testes do projeto (Preparando o terreno para o futuro)
test:
	@echo "Rodando testes..."
	go test ./...