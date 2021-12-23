CREATE TABLE public."ClientConfig"
(
    "id"            bigint                      NOT NULL,
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

CREATE TABLE public."Application"
(
    "id"         bigint                      NOT NULL,
    "createdAt"  timestamp without time zone NOT NULL,
    "updatedAt"  timestamp without time zone NOT NULL,
    "metadata"   JSON                        NULL,
    "modifiedBy" JSON                        NOT NULL,
    "version"    bigint                      NOT NULL DEFAULT 1,
    "name"       varchar(255)                NOT NULL,
    "xid"        varchar(64)                 NOT NULL,
    PRIMARY KEY ("id")
);

ALTER TABLE public."ClientConfig"
    ADD CONSTRAINT "FK_Application__applicationId" FOREIGN KEY ("applicationId") REFERENCES public."Application" (id);