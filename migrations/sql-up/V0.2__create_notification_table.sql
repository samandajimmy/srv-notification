-- Create Tables
CREATE TABLE "public"."Notification"
(
    "id"            uuid                        NOT NULL,
    "createdAt"     timestamp without time zone NOT NULL,
    "updatedAt"     timestamp without time zone NULL,
    "modifiedBy"    json                        NOT NULL,
    "metadata"      json                        NULL,
    "version"       bigint                      NOT NULL DEFAULT 1,
    "userRefId"     bigint,
    "title"         varchar(255)                NOT NULL,
    "tagline"       varchar(255),
    "content_short" varchar(255)                NULL,
    "content"       TEXT                        NULL,
    "isRead"        bool                                 DEFAULT FALSE,
    "readAt"        timestamp without time zone NULL,
    PRIMARY KEY ("id")
);

-- Create Index for table "Notification"
CREATE INDEX ON public."Notification" ("userRefId");
