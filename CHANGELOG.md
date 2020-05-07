# 2020-05-08

  * Added support for using a secondary calendar.  Use of this feature requires
    invalidating all previous user sessions in deployed instances since the
    required OAuth2 scope has changed.  This can be accomplished either by
    removing the entries from the datastore or changing the session secret.

# 2020-04-17

  * Initial release.
