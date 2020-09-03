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
  * (If EHR writeback is enabled) the EHR address and port to receive HL7 messages sent with MLLP

The OAuth2 Client ID can be created using the [Cloud
Console](https://console.cloud.google.com/apis/credentials/consent).
  * Select 'Web application' as the type. 
  * The redirect URI should point to '/authenticate' on the appropriate
    server (e.g., localhost or your appspot.com subdomain).
  * In the "Scopes for Google APIs" section of the OAuth2 consent screen, make
    sure you add `https://www.googleapis.com/auth/calendar.events` as a scope.

To provide these settings, create a file called `settings.json` using the
instructions in `settings.json-example`.

## Choosing the calendar for events

The example settings file creates calendar events on the primary calendar.
Instead, a custom name can be specified and the named calendar will be created
(if it doesn't exist) and events will be added there.

*Note that if you choose to use this feature you must invalidate all previous
user sessions since the required OAuth2 scopes are broader.*

## SMART on FHIR configuration

The launch URL should be set to `/launch.html` on the appropriate server (e.g.
localhost or your appspot.com subdomain).

Note that this application does not work inside a frame, so it must be
configured to launch as a new window in the SMART on FHIR integration point.

## EHR writeback configuration

EHR writeback can be enabled/disabled in settings.json.
If enabled, the app will report patient arrived events and appointment status changes
by sending MLLP-encoded HL7 messages to the specified EHR address and port.
The running environment must be allowed to establish TCP socket with the EHR server.

# Running locally

You will need a recent version of Node.  Once installed, you can run `npm
install` to install the other required dependencies.

Once everything is installed, configure Google Application Default credentials
with access to a Cloud Datastore in a project you own and run `npm start`.

# Deploying on Google Cloud

To deploy on Google Cloud, you will need a project that does not already have
an Appengine application deployed.

Deploy the application using `gcloud app deploy`.

If EHR writeback is enabled, a serverless connector must be created for the app
to access resources in a VPC network where the EHR server is connected to. Then
to the [app.yaml](https://github.com/google/meet-on-fhir/blob/5c71b37b3bdf0703c281bd8e23d5dd383b28bee8/app.yaml)
file in the root directory of this repository, add the following section:

```
vpc_access_connector:
    name: projects/PROJECT_ID/locations/CONNECTOR_REGION/connectors/CONNECTOR_NAME
```
Refer to (https://cloud.google.com/appengine/docs/standard/nodejs/connecting-vpc) for
more details.

# Testing

You can use the SMART on FHIR [launcher](https://launch.smarthealthit.org/) to
make sure your instance is working as expected.  Open two different profiles in
a web browser and configure one as a patient and one as a physician and click
the launch button.
