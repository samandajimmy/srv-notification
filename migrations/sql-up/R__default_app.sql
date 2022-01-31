INSERT INTO public."Application" ("createdAt", "updatedAt", "modifiedBy", version, metadata, name, xid)
VALUES ('2020-01-01 00:00:00', '2020-01-01 00:00:00', '{"id": "0","role": "SEEDER","fullName": "Seeder"}', 1, '{}', 'Default Configuration', 'DEFAULT_CONFIG')
ON CONFLICT DO NOTHING;
