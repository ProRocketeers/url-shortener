-- Atlas requires a clean, empty database in order to create migrations, like Prisma's shadow database
-- this script is mounted in the local development Postgres container and executed on container create
CREATE DATABASE dev;