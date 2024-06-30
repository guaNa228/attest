# Attestation Service Backend

This repository contains the backend code for an attestation service used in a university setting. The service is developed in Go and provides various functionalities related to student attestation and administrative operations.

## Features

- **User Authentication**: Secure login system for users.
- **Attestation Management**: Handle student attestation processes efficiently.
- **Class and Group Management**: Manage classes and groups within the system.
- **Email Notifications**: Automated email sending for various notifications.
- **Readiness Check**: Endpoint for checking the service status.
- **Database Management**: SQL scripts and handlers for database operations.
- **Automatic Web Scraping**: Automated scraping for group-class-teacher management.

## Project Structure

- `internal`: Internal application logic.
- `logger`: Logging functionalities.
- `parsing`: Data parsing utilities.
- `sql`: SQL scripts and database interaction code.
- `translit`: Transliteration utilities.
- `vendor`: External dependencies.
- Various Go files implementing specific handlers for user, classes, attestation, emails, and more.

## Getting Started

### Prerequisites

- Go 1.16 or higher
- PostgreSQL

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/guaNa228/attest.git
   cd attest
