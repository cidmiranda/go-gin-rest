# Go RESTful API

Aplicação rest, desenvolvida em [Go](https://go.dev/), [Fiber](https://gofiber.io/) and [PostgreSQL](https://www.postgresql.org/).

## Executando a aplicação

1. Instalar o [Go](https://go.dev/dl/)
2. Instalar o [Docker Compose](https://docs.docker.com/compose/install/). Ele precisa ser iniciado.
3. Clonar a aplicação:

```bash
git clone git@github.com:cidmiranda/go-fiber-rest.git
```

4. Vá até a pasta da aplicação

```bash
cd go-fiber-rest
```

5. Execute a imagem docker

```bash
docker-compose -f docker-compose.yml up --build
```

6. Acesse a documentação da aplicação em http://localhost:8080/swagger/index.html.
