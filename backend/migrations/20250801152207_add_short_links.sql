-- Create "short_links" table
CREATE TABLE "public"."short_links" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "original_url" text NOT NULL,
  "slug" text NOT NULL,
  "expires_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_short_links_slug" UNIQUE ("slug")
);
-- Create index "idx_short_links_deleted_at" to table: "short_links"
CREATE INDEX "idx_short_links_deleted_at" ON "public"."short_links" ("deleted_at");
