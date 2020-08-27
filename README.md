# TeleHealth v1

# Running locally

Go to the goserver directory and run
`go run main.go --authorized_fhir_url={fhir_base_url}`

# Testing

You can use the SMART on FHIR [launcher](https://launch.smarthealthit.org/) to
make sure your instance is working as expected.  Open two different profiles in
a web browser and configure one as a patient and one as a physician and click
the launch button.

If testing locally, first run the server with
`--authorized_fhir_url=https://r4.smarthealthit.org`
then use "http://localhost:8080" as the App Launch URL. 
