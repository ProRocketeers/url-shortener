data "external_schema" "gorm" {
  // the command that creates an SQL schema and passes it to Atlas
  program = [
    "go", "run", "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./domain/model", 
    "--dialect", "postgres"
  ]
  // Atlas automatically reads the structs that embed the `gorm.Model` struct or use `gorm` tags and constructs a database schema from it
}

env "local" {
  // this references the source of the SQL schema. Why does Atlas use `url` property for it is beyond my fuckin imagination
  src = data.external_schema.gorm.url
  // local dev DB that's used for creating migrations - see `docker-compose.yaml`
  dev = "postgresql://postgres:postgres@localhost:5432/dev?sslmode=disable&timezone=UTC"
  url = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable&timezone=UTC"

  migration {
    dir = "file://domain/migrations"
  }

  format {
    migrate {
      // this tells Atlas how to format the SQL migrations - purely for readability purposes
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "runtime" {
  url = getenv("ATLAS_DATABASE_URL")
}