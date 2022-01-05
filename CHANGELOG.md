# Changelog

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