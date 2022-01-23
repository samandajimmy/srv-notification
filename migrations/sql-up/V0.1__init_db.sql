-- Create tables

CREATE TABLE public."Application"
(
    "id"         bigserial                   NOT NULL,
    "createdAt"  timestamp without time zone NOT NULL,
    "updatedAt"  timestamp without time zone NOT NULL,
    "metadata"   JSON                        NULL,
    "modifiedBy" JSON                        NOT NULL,
    "version"    bigint                      NOT NULL DEFAULT 1,
    "name"       varchar(255)                NOT NULL,
    "xid"        varchar(64)                 NOT NULL,
    "apiKey"     varchar(255)                NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE public."ClientConfig"
(
    "id"            bigserial                   NOT NULL,
    "createdAt"     timestamp without time zone NOT NULL,
    "updatedAt"     timestamp without time zone NOT NULL,
    "metadata"      JSON                        NULL,
    "modifiedBy"    JSON                        NOT NULL,
    "version"       bigint                      NOT NULL DEFAULT 1,
    "key"           varchar(255)                NOT NULL,
    "value"         JSON                        NOT NULL,
    "applicationId" smallint                    NOT NULL,
    "xid"           varchar(64)                 NOT NULL,
    PRIMARY KEY ("id")
);


-- Create Index for table "Application"

CREATE UNIQUE INDEX ON public."Application" (xid);
CREATE INDEX ON public."Application" (name);

-- Create Index for table "ClientConfig"

ALTER TABLE public."ClientConfig"
    ADD CONSTRAINT "FK_Application__applicationId" FOREIGN KEY ("applicationId") REFERENCES public."Application" (id);

CREATE INDEX ON public."ClientConfig" (xid);
CREATE INDEX ON public."ClientConfig" ("applicationId");
CREATE INDEX ON public."ClientConfig" ("key");