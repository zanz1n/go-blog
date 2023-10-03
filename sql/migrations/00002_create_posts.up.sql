CREATE TABLE "posts" (
    "id" uuid PRIMARY KEY,
    "created_at" timestamp(3) NOT NULL DEFAULT current_timestamp,
    "updated_at" timestamp(3) NOT NULL DEFAULT current_timestamp,
    "title" varchar(192) NOT NULL,
    "content" bytea NOT NULL,
    "thumb_image" uuid,
    "user_id" uuid NOT NULL,

    CONSTRAINT "posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "user"("id")
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX "posts_user_id_idx" ON "posts"("user_id");
CREATE INDEX "posts_created_at_idx" ON "posts"("created_at");
