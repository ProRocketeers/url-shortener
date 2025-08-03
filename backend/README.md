# 🚀 URL Shortener backend
This is a backend service for the URL shortener project. It provides a small REST API for shortening URLs and following those shortened links

## How to run
1. First, install `Go` - you need at least version `1.24.5`, verify with `go version`

2. Then install dependencies:
    - `swaggo` - for generating API docs
      ```bash
      go install github.com/swaggo/swag/cmd/swag@latest
      swag --version
      ```
    - `air` - for hot-reloading during development
      ```bash
      go install github.com/air-verse/air@latest
      air -v
      ```
    - `atlas` - for handling database migrations
      ```bash
      curl -sSf https://atlasgo.sh | sh
      go get ariga.io/atlas-provider-gorm
      ```

3. Copy the `.env.example` as `.env` and edit as necessary
    ```bash
    cp .env.example .env
    ```
4. Create and migrate the database
    ```bash
    docker compose up -d
    make migrate-dev
    ```
5. Start the server with hot-reloading
    ```bash
    make run
    ```

### API docs
To regenerate docs from Swagger annotations/comments, run
```bash
make gen-swagger
```

### Migrations
To create a new migration after updating the models, run
```bash
make create-migration NAME=migration_name
```

And then apply to your local DB 
```bash
make migrate-dev
```