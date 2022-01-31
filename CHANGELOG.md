# Changelog

## 0.9.0

- [CHANGED] Return notification detail on Create Notification
- [ADDED] API Get Detail Notification

## 0.8.1

- [CHANGED] Add prefix for migration config

## 0.8.0

- [FIXED] set apiKey when create update delete application
- [ADDED] Add Column apiKey to Application table
- [ADDED] implement create notification when send push notification
- [ADDED] Add apiKey to Application table and implement apikey when send notification
- [ADDED] Send notification handler and pubsub
- [ADDED] Create notification service
- [ADDED] notification repository
- [ADDED] notification table schema

## 0.7.1

- [ADDED] Refactor client config detail response and added filter by application xid

## 0.7.0

- [ADDED] Update client config
- [CHANGED] Modifier `FullName` using camelCase from snake case
- [ADDED] Route delete and detail client config
- [ADDED] List client config

## 0.6.0

- [FIXED] Remove formatting app name on create
- [ADDED] Update application service.
- [ADDED] List application service.
- [ADDED] Delete application service.
- [FIXED] Generalize application name to uppercase.
- [ADDED] Get detail application service
- [ADDED] Initialize handler for application service
- [ADDED] Create application service.
- [ADDED] Implement load config from database send fcm notification.
- [FIXED] Remove smtp and firebase config from env
- [ADDED] Implement load config from database send email.
- [ADDED] Implement load config from database send fcm notification.
- [ADDED] Initialize config from database.
- [ADDED] Initialize migration for multiple client management.

## 0.5.1

- [FIXED] Remove payload message on sending email.

## 0.5.0

- [CHANGED] Implement context aware Service
- [CHANGED] Fallback to logger context if set
- [ADDED] Add context aware logger and implement in on create log child
- [FIXED] Change email attachment path
- [FIXED] Remove config loader on Core boot

## 0.4.1

- [FIXED] Fix import

## 0.4.0

- [ADDED] Implement set request id on request context
- [CHANGED] Implement json formatted logging
- [CHANGED] Replace logging dependency with nbs-go/nlogger

## 0.3.5

- [FIXED] Prevent retrying pubsub on error sending email
- [FIXED] Prevent retrying pubsub on error sending push notification

## 0.3.4

- [CHANGED] Remove firebase credential from build image
- [CHANGED] Move Firebase Service Account from file to env
- [FIXED] Add to email format validation

## 0.3.3

- [FIXED] Prevent lookup for MX Record validation for From email

## 0.3.2

- [FIXED] Fix nested struct validation

## 0.3.1

- [FIXED] Add logging on push notification sent
- [FIXED] Publish event on fcm request

## 0.3.0

- [ADDED] Implement async send fcm push notification

## 0.2.0

- [ADDED] Implement send email async with PubSub

## 0.1.6

- [FIXED] Missing tmp dir
- [FIXED] Missing /app dir

## 0.1.5

- [FIXED] Fix missing responses.yml and firebase-secret.json file on build docker

## 0.1.4

- [FIXED] Copy firebase-secret on building release container

## 0.1.3

- [FIXED] Add log-in to artifactory for pull docker images

## 0.1.2

- [CHANGED] Change docker image namespace to Pegadaian Artifactory
- [FIXED] Skip database initialization
- [FIXED] Fix build script
- [FIXED] Add firebase secret on build
- [FIXED] Install make on build stage
- [CHANGED] Add port env on build
- [CHANGED] Rename package to match repository

## 0.1.1

- [CHANGED] Add port env on build
- [CHANGED] Rename package to match repository

## 0.1.0

- [CHANGED] Rename package to match repository
- [ADDED] Dynamic attribute data payload notification
- [FIXED] Fix bugs send email without attachment handling
- [FIXED] add log handling when validate payload notification
- [FIXED] fix call validate for payload notification
- [ADDED] add endpoint send notification FCM API
- [FIXED] add log handling when validate payload notification
- [FIXED] fix call validate for payload notification
- [CHANGED] move send email logic to pkg nmail.
- [FIXED] fix dto validate from email isEmail.
- [CHANGED] move send email logic to pkg nmail.
- [FIXED] fix dto validate from email isEmail.
- [FIXED] fix add payload validation on handler
- [FIXED] fix add payload validation on handler
- [FIXED] fix validation for email must be valid email
- [FIXED] fix validation for email must be valid email
- [FIXED] payload dto for send email add name and email.
- [FIXED] route name send email to singular.
- [FIXED] fixed add attribute dto for send email.
- [FIXED] payload dto for send email add name and email.
- [FIXED] route name send email to singular.
- [FIXED] fixed add attribute dto for send email.
- [ADDED] add endpoint send notification FCM API
- [FIXED] Fixed typo service name
- [ADDED] add pkg firebase for FCM
- [ADDED] add endpoint Send Email API
- [ADDED] Env for smtp and add dependencies.