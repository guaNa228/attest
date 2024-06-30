Attestation Service Backend

This repository contains the backend code for an attestation service used in a university setting. The service is developed in Go and provides various functionalities related to student attestation and administrative operations.
Features

    User Authentication: Secure login system for users.
    Attestation Management: Handle student attestation processes efficiently.
    Class and Group Management: Manage classes and groups within the system.
    Email Notifications: Automated email sending for various notifications.
    Readiness Check: Endpoint for checking the service status.
    Database Management: SQL scripts and handlers for database operations.

Project Structure

    internal: Internal application logic.
    logger: Logging functionalities.
    parsing: Data parsing utilities.
    sql: SQL scripts and database interaction code.
    translit: Transliteration utilities.
    vendor: External dependencies.
    Various Go files implementing specific handlers for user, classes, attestation, emails, and more.

Getting Started
Prerequisites

    Go 1.16 or higher
    PostgreSQL

Installation

    Clone the repository:

    sh

git clone https://github.com/guaNa228/attest.git
cd attest

Install dependencies:

sh

    go mod download

    Set up the database using the scripts in the sql directory.

Running the Service

    Compile the code:

    sh

go build -o attest

Run the service:

sh

    ./attest

Usage

    Authentication: Secure endpoints using the provided middleware.
    Endpoints: Various endpoints for managing students, classes, groups, and attestation records.
    Email Notifications: Automatically send emails upon specific actions.

Contribution

    Fork the repository.
    Create a new branch (git checkout -b feature-branch).
    Make your changes.
    Commit your changes (git commit -m 'Add new feature').
    Push to the branch (git push origin feature-branch).
    Create a Pull Request.
