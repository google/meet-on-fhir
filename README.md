# TeleHealth Integration Reference Implementation

This repository contains an early release reference implementation that may be
useful to developers wishing to integrate Google Meet into EHR systems that
support SMART on FHIR.

To use this reference implementation, you should fork this repository and add
any required user interface or other integrations that are appropriate.

## Prerequisites

This reference implementation assumes that:

  * The provider/physician has access to a GSuite account with the Calendar and
    Meets applications enabled.
  * If the patient is accessing the application on the desktop, they can join
    the meeting anonymously.  If the patient accesses the application from
    a mobile device, they must download the Meet application and sign in.

# Configuration

The application requires several things to be configured:

  * An OAuth2 Client ID and secret to access the Google Cloud Datastore
  * A SMART on FHIR Client ID registered with the EHR system
  * A secret key used to encrypt the session cookie
  * The calendar API must be enabled in the project

The OAuth2 Client ID can be created using the [Cloud
Console](https://console.cloud.google.com/apis/credentials/consent).
  * Select 'Web application' as the type. 
  * The redirect URI should point to '/authenticate' on the appropriate
    server (e.g., localhost or your appspot.com subdomain).
  * In the "Scopes for Google APIs" section of the OAuth2 consent screen, make
    sure you add `https://www.googleapis.com/auth/calendar.events` as a scope.

To provide these settings, create a file called `settings.json` using the
instructions in `settings.json-example`.

## SMART on FHIR configuration

The launch URL should be set to `/launch.html` on the appropriate server (e.g.
localhost or your appspot.com subdomain).

Note that this application does not work inside a frame, so it must be
configured to launch as a new window in the SMART on FHIR integration point.

# Runing locally

Configure Google Application Default credentials and run `npm start`.

# Deploying on Google Cloud

To deploy on Google Cloud, you will need a project that does not already have
an Appengine application deployed.

Deploy the application using `gcloud app deploy`.

# Testing

You can use the SMART on FHIR [launcher](https://launch.smarthealthit.org/) to
make sure your instance is working as expected.  Open two different profiles in
a web browser and configure one as a patient and one as a physician and click
the launch button.
