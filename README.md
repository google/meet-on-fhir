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

# Running locally

Go to the goserver directory and run

`go run main.go --authorized_fhir_url={fhir_base_url}`

The --authorized_fhir_url parameter must be set so that the app is authorized to launch.
The value must match the iss query parameter passed to the server when calling its launch endpoint.

# Testing

You can use the SMART on FHIR [launcher](https://launch.smarthealthit.org/) to
make sure your instance is working as expected.  Open two different profiles in
a web browser and configure one as a patient and one as a physician and click
the launch button.

If testing locally, first run the server with
`--authorized_fhir_url=https://r4.smarthealthit.org`
then use "http://localhost:8080" as the App Launch URL. 

# UI customization

## Change texts (such as patient consent message) in the UI
All customizable texts are saved in the ui/src/i18n-strings.ts file. If you want to replace a
displayed text in UI, find the text in the file first and replace it with your preferred text. If
a field has an array of texts, each text is displayed as a paragraph in the UI.

## Add more supported languages
We only provide English support out of box. However, we have implemented a simple framework which
allows you to easily add more languages. Not all languages are supported at the moment (e.g RTL
languages). Please refer to the LanguageCode enum in ui/src/i18n-helper.ts as the source of truth.
If you can't find a language in the enum, it means you can't add the language yet.

To add a new language, please follow these steps:
1. Find the language in LanguageCode and remember the 2-4 letter code (e.g. 'es' for Spanish)
corresponding to the language. 
1. Open ui/src/i18n-strings.ts file.
1. Copy and paste the English entry ('en') and then change 'en' to the 2-4 letters from above.
1. Replace all strings in the new entry with the correct translated strings.

Once one or more languages are added, users will be presented with a language selector initially.
