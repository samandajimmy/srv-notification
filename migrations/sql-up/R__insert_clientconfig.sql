INSERT INTO public."ClientConfig"(id, "createdAt", "updatedAt", "modifiedBy", version, metadata, key, value,
                                  "applicationId", "xid")
VALUES (1, '2020-01-01 00:00:00', '2020-01-01 00:00:00', '{
  "id": "0",
  "role": "SEEDER",
  "fullName": "Seeder"
}', 1, '{}',
        'FIREBASE_SERVICE_ACCOUNT_CRED',
        '{
          "type": "service_account",
          "project_id": "pds-dev-f65b3",
          "private_key_id": "a56972c48ae8366764d3215d045d4978f45b3da1",
          "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDAfFlgAfICB7XJ\ndRITIUPYB5cqVOvXmZSlml1eVa967dZgnn2bVf3gWI4p2hi3bsXAz8a9TCNw6j9p\nsU2Mp1FlmcjDlnBDAdoLIiD6MzGYBag5k8Cw4eXayaNxPq5mO0ddHdGa/tqvdQYH\nbTx8W/AR3X5vcvzKynZUdGME3KadteUo3iyfTeXro3KVkX9Y9XPafnQHWUjuBY0h\nIr8WHjRbVpEeuFL2WsZHEzEhO0UIdKsOb22Dr3iFTjJefsHVaaroh8ZvTVfk+cCB\nRgyWDDlAaJMOr1OcB+hn0rlzRnOVmdCv3jrk4RLMtRVxFOnzc53VedVCOo11EFi8\nVJR+n9EnAgMBAAECggEAULFQfnESVUuKJ+ROKXreiCfWdUaYgA/AQxRNZAijwfMG\n1gZbPc102qIFJwJpLizf9g6kkCLlEKcC4noTuo1CEEfB5Eyiz1RtZhFupbTka4ij\nl0+bjguFYoz5WbYfQnhDWIPxpPqGDtwOJPrkSnX4VGT4ZhxcYV8y8ADCqf0eVCp8\nX1wdzPApsxLSOr3ig9SyETU5OddnC/+acBT/vsjLL4pLA8AWcUAgMlsvvY5bWscL\nod6v9FZP0cgs4X1SYba1WQPTq6bAfE1eRJYxWmBdBjL5FDnlOErBTl9QIQc1R+hr\n5v64E4ZaI7qdXEdl+NAy4DLbpnySyoesTUk1788eMQKBgQDsxTyGj1lQBT7Wupci\nwOyH4k9TVdZyM2l1GmOk1Y680n0B3StwWD2v0Vl5VrDYLoPyC7PM0MA2kC1X93QH\nMikOMn0rU4xv64VCexzLqQvkMljknkgM3eXcZaNDATjFTQ1sHeh8pVR4qiuT62Wt\nLjcO35xsHRB2afW7pfgg3iQPYwKBgQDQHmBRw0T5daQ86PjvBJ1IRgOI1sberGJz\nAl9Q7oCArn+/nZUQhhY3Xb8277MZQlpdgBfNUnkkeY3sHBF1MOleOlN3KjE90dr3\nGe3vuWuH7Oh/PfdC/sMesJLOMsW+ZJkar0Ssp12DNk19D1cRrrkRLwX7dgEin+NL\neOSmO3LsbQKBgAlV3YUQsdzN2CRvRvY/1ROmgKowgDwQet/7ImKlaPNY+UTRi5zq\nXcRI5NY77M0ZSGqKu5QfxvRfyunk/9YozCWbKARFTww6pQ5x/DireaSNt6OL+htH\nxIBkIYPK0Io294iDxV7kxefcDcvPRDsHz3PurSQ2ISgKFX5IlPf2ykUxAoGBAI/Q\nHdECNZyIXYi1mKdaMfFqaDDb8aqXxqQgxIrhdKz3aFGZ7BLyBVIXFvY4ZzOSNW2d\nAVWhoxLAaID62Fl6BhlWBq0227YTWNMd+NyJ7bOM0xByWnXSJDUF4TxZu2mYjG/z\nI2qHcMgl8x/zRMB0U3B7ZQ/h+GDbya4yiRYRyJV1AoGBAJ5ucY1CPv9+QgmQZJ5X\n5gi/4jD6UN/WkCP7pNgD7Pl1EVCnIS/lc/7LoK8dDjrV9NfhBw3+qkiBAL0mM3RW\nwwhGIaM/DvzpvHQ5umxR2w5V3lFelF0hCLpCndpdVZIYv0plwXE7JObGoHyCTTq4\nlWZaPacks1knhNRdy5UilRFO\n-----END PRIVATE KEY-----\n",
          "client_email": "firebase-adminsdk-dr4fd@pds-dev-f65b3.iam.gserviceaccount.com",
          "client_id": "118349119665464491067",
          "auth_uri": "https://accounts.google.com/o/oauth2/auth",
          "token_uri": "https://oauth2.googleapis.com/token",
          "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
          "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-dr4fd%40pds-dev-f65b3.iam.gserviceaccount.com"
        }',
        1, 'PDSFRBS'),
       (2, '2020-01-01 00:00:00', '2020-01-01 00:00:00', '{
         "id": "0",
         "role": "SEEDER",
         "fullName": "Seeder"
       }', 1, '{}',
        'SMTP',
        '{
          "SMTP_HOST": "smtp.eu.mailgun.org",
          "SMTP_PORT": "465",
          "SMTP_USERNAME": "dev_pds@mg.nbs.co.id",
          "SMTP_PASSWORD": "e738b1fb441633906274b63ca515827e-10eedde5-288e6771"
        }',
        1, 'PDSSMTP'),
       (3, '2020-01-01 00:00:00', '2020-01-01 00:00:00', '{
         "id": "0",
         "role": "SEEDER",
         "fullName": "Seeder"
       }', 1, '{}',
        'FIREBASE_SERVICE_ACCOUNT_CRED',
        '{
          "type": "service_account",
          "project_id": "pds-dev-f65b3",
          "private_key_id": "a56972c48ae8366764d3215d045d4978f45b3da1",
          "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDAfFlgAfICB7XJ\ndRITIUPYB5cqVOvXmZSlml1eVa967dZgnn2bVf3gWI4p2hi3bsXAz8a9TCNw6j9p\nsU2Mp1FlmcjDlnBDAdoLIiD6MzGYBag5k8Cw4eXayaNxPq5mO0ddHdGa/tqvdQYH\nbTx8W/AR3X5vcvzKynZUdGME3KadteUo3iyfTeXro3KVkX9Y9XPafnQHWUjuBY0h\nIr8WHjRbVpEeuFL2WsZHEzEhO0UIdKsOb22Dr3iFTjJefsHVaaroh8ZvTVfk+cCB\nRgyWDDlAaJMOr1OcB+hn0rlzRnOVmdCv3jrk4RLMtRVxFOnzc53VedVCOo11EFi8\nVJR+n9EnAgMBAAECggEAULFQfnESVUuKJ+ROKXreiCfWdUaYgA/AQxRNZAijwfMG\n1gZbPc102qIFJwJpLizf9g6kkCLlEKcC4noTuo1CEEfB5Eyiz1RtZhFupbTka4ij\nl0+bjguFYoz5WbYfQnhDWIPxpPqGDtwOJPrkSnX4VGT4ZhxcYV8y8ADCqf0eVCp8\nX1wdzPApsxLSOr3ig9SyETU5OddnC/+acBT/vsjLL4pLA8AWcUAgMlsvvY5bWscL\nod6v9FZP0cgs4X1SYba1WQPTq6bAfE1eRJYxWmBdBjL5FDnlOErBTl9QIQc1R+hr\n5v64E4ZaI7qdXEdl+NAy4DLbpnySyoesTUk1788eMQKBgQDsxTyGj1lQBT7Wupci\nwOyH4k9TVdZyM2l1GmOk1Y680n0B3StwWD2v0Vl5VrDYLoPyC7PM0MA2kC1X93QH\nMikOMn0rU4xv64VCexzLqQvkMljknkgM3eXcZaNDATjFTQ1sHeh8pVR4qiuT62Wt\nLjcO35xsHRB2afW7pfgg3iQPYwKBgQDQHmBRw0T5daQ86PjvBJ1IRgOI1sberGJz\nAl9Q7oCArn+/nZUQhhY3Xb8277MZQlpdgBfNUnkkeY3sHBF1MOleOlN3KjE90dr3\nGe3vuWuH7Oh/PfdC/sMesJLOMsW+ZJkar0Ssp12DNk19D1cRrrkRLwX7dgEin+NL\neOSmO3LsbQKBgAlV3YUQsdzN2CRvRvY/1ROmgKowgDwQet/7ImKlaPNY+UTRi5zq\nXcRI5NY77M0ZSGqKu5QfxvRfyunk/9YozCWbKARFTww6pQ5x/DireaSNt6OL+htH\nxIBkIYPK0Io294iDxV7kxefcDcvPRDsHz3PurSQ2ISgKFX5IlPf2ykUxAoGBAI/Q\nHdECNZyIXYi1mKdaMfFqaDDb8aqXxqQgxIrhdKz3aFGZ7BLyBVIXFvY4ZzOSNW2d\nAVWhoxLAaID62Fl6BhlWBq0227YTWNMd+NyJ7bOM0xByWnXSJDUF4TxZu2mYjG/z\nI2qHcMgl8x/zRMB0U3B7ZQ/h+GDbya4yiRYRyJV1AoGBAJ5ucY1CPv9+QgmQZJ5X\n5gi/4jD6UN/WkCP7pNgD7Pl1EVCnIS/lc/7LoK8dDjrV9NfhBw3+qkiBAL0mM3RW\nwwhGIaM/DvzpvHQ5umxR2w5V3lFelF0hCLpCndpdVZIYv0plwXE7JObGoHyCTTq4\nlWZaPacks1knhNRdy5UilRFO\n-----END PRIVATE KEY-----\n",
          "client_email": "firebase-adminsdk-dr4fd@pds-dev-f65b3.iam.gserviceaccount.com",
          "client_id": "118349119665464491067",
          "auth_uri": "https://accounts.google.com/o/oauth2/auth",
          "token_uri": "https://oauth2.googleapis.com/token",
          "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
          "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-dr4fd%40pds-dev-f65b3.iam.gserviceaccount.com"
        }',
        2, 'PSDSFRBS'),
       (4, '2020-01-01 00:00:00', '2020-01-01 00:00:00', '{
         "id": "0",
         "role": "SEEDER",
         "fullName": "Seeder"
       }', 1, '{}',
        'SMTP',
        '{
          "SMTP_HOST": "smtp.eu.mailgun.org",
          "SMTP_PORT": "465",
          "SMTP_USERNAME": "dev_pds@mg.nbs.co.id",
          "SMTP_PASSWORD": "e738b1fb441633906274b63ca515827e-10eedde5-288e6771"
        }',
        2, 'PSDSSMTP')
ON CONFLICT DO NOTHING;
