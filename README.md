# Kuduga

Kuduga is an Indonesian digital-product purchasing app for mobile credit, data packages, e-wallet top-ups, electricity tokens, game top-ups, and bill payments. This repository contains its Go backend and Flutter mobile app.

## Project story

I started Kuduga in **August 2023** as my first real-world Flutter project. I later shut it down due to marketing challenges, funding issues, and not attracting enough users, and decided to share the source code.

## Features

- Digital-product purchases and prepaid/postpaid bill payments.
- Account balances, deposits, transfers, and transaction history.
- Email OTP and Google sign-in, PIN settings, and biometric app locking.
- Push notifications, referrals, identity verification, and support tickets.
- Purchase statistics and thermal receipt printing.

## Screenshots

<p align="center">
  <img src="https://github.com/user-attachments/assets/360d07cf-7220-43ee-8148-abe3f8122eb6" alt="Kuduga app overview" width="400">
</p>

<p align="center">
  <img src="https://github.com/user-attachments/assets/3be673f4-ce52-4e24-baac-8bc4ca605907" alt="Kuduga app screenshot 1" height="400">
  <img src="https://github.com/user-attachments/assets/933cdbdf-e319-44c5-828e-753fc9791d4a" alt="Kuduga app screenshot 2" height="400">
  <img src="https://github.com/user-attachments/assets/829498a9-6643-41f1-a3be-772562d3852f" alt="Kuduga app screenshot 3" height="400">
  <img src="https://github.com/user-attachments/assets/b17bab87-7be7-4ada-839f-97c059a767f0" alt="Kuduga app screenshot 4" height="400">
  <img src="https://github.com/user-attachments/assets/b226385d-dae2-4e93-907e-e12e419b156b" alt="Kuduga app screenshot 5" height="400">
  <img src="https://github.com/user-attachments/assets/f99a38a2-bd28-498a-b0b7-83ad72c40cfc" alt="Kuduga app screenshot 6" height="400">
</p>

## Technology

The backend uses Go, Fiber, MySQL, and Redis. The mobile app uses Flutter and Dart, with Android and iOS projects. Integrations include Firebase, Google sign-in, MinIO, AWS SES, Digiflazz, Duitku, and a BCA payment service.

## Repository structure

- [`finance/`](finance/) — Go backend, database migrations, and integrations.
- [`finance_app/`](finance_app/) — Flutter mobile app and assets.

## Getting started

Use your own services and credentials. Configuration examples are included.

### Backend

Install Go 1.21 or later and set up MySQL, Redis, and MinIO. Copy the environment example:

```sh
cd finance
cp .env-example .env
```

Fill in `.env`, including your Firebase Admin credentials via `GOOGLE_APPLICATION_CREDENTIALS`, then start the backend:

```sh
go run .
```

Migrations create the database schema. Add your application metadata and products, and configure email and payment integrations as needed.

### Mobile app

Install Flutter, then:

```sh
cd finance_app
flutter pub get
```

Configure your own Firebase project using `flutterfire configure`, including Android/iOS configuration and Google sign-in. Then run the app against your backend. For an Android emulator:

```sh
flutter run --dart-define=API_HOST=10.0.2.2:3001 --dart-define=API_SECURE=false
```

For a physical device, set `API_HOST` to your backend's reachable host and port. Use `API_SECURE=true` for HTTPS.

Replace `example.invalid` placeholders with your own links and contact details.
