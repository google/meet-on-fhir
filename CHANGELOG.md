# 2020-05-19

  * Added example UI skinning of Waiting Room.
  * Added Consent screen for patients required before entering Waiting Room.
  * Added multi-language support for English and Spanish UIs.
  * Added fallback support for older Epic FHIR servers that can't pass back an id_token and can only pass extra CONTEXT information outside of the id_token JWT.
  * Fixed clearing the waitFor method timer after successful return of Meet URL so that mobile website invocations won't continue to execute timer because they switch to the App Store to open the Google Meet application rather than redirect to the Meet URL in the browser.
  * Added cloudbuild.yaml and build steps for a Cloud Build trigger.

# 2020-05-08

  * Added support for using a secondary calendar.  Use of this feature requires
    invalidating all previous user sessions in deployed instances since the
    required OAuth2 scope has changed.  This can be accomplished either by
    removing the entries from the datastore or changing the session secret.

# 2020-04-17

  * Initial release.
