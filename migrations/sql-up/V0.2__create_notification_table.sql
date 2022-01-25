-- Create Tables
CREATE TABLE "public"."Notification"
(
    "id"             uuid                        NOT NULL,
    "createdAt"      timestamp without time zone NOT NULL,
    "updatedAt"      timestamp without time zone NULL,
    "modifiedBy"     json                        NOT NULL,
    "metadata"       json                        NULL,
    "version"        bigint                      NOT NULL DEFAULT 1,
    "applicationId"  bigint                      NOT NULL,
    "userRefId"      bigint,
    "isRead"         bool                        NULL     DEFAULT FALSE,
    "readAt"         timestamp without time zone NULL,
    "options"        json                        NULL,
    PRIMARY KEY ("id")
);

-- Create Index for table "Notification"
CREATE INDEX ON public."Notification" ("userRefId");
CREATE INDEX ON public."Notification" ("applicationId");


-- Alter Application
alter table "Application"
    add "apiKey" varchar(255);
