-- Create "request_infos" table
CREATE TABLE "public"."request_infos" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "request_id" text NOT NULL,
  "timestamp" timestamptz NOT NULL,
  "real_ip" text NULL,
  "user_agent" text NULL,
  "headers" jsonb NULL,
  "path" text NOT NULL,
  "method" text NOT NULL,
  "query" jsonb NULL,
  "body" jsonb NULL,
  PRIMARY KEY ("id")
);
