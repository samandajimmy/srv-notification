INSERT INTO public."Application"(id, "createdAt", "updatedAt", "modifiedBy", version, metadata, name, xid)
VALUES (1, '2020-01-01 00:00:00', '2020-01-01 00:00:00', '{"id": "0","role": "SEEDER","fullName": "Seeder"}', 1, '{}', 'PDS', 'APKPDS'),
       (2, '2020-01-01 00:00:00', '2020-01-01 00:00:00', '{"id": "0","role": "SEEDER","fullName": "Seeder"}', 1, '{}', 'PSDS', 'APKPSDS')
ON CONFLICT DO NOTHING;
