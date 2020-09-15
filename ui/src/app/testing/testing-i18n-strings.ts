import {ALL_LANGUAGES_ASSETS_DATA} from '../i18n-strings';

export const TESTING_ALL_LANGUAGE_ASSETS_DATA = {
  'en': {submitButton: 'English',},
  'es': {submitButton: 'Spanish'},
};

/**
 * This provider provides ALL_LANGUAGE_ASSETS_DATA for testing.
 */
export const TESTING_ALL_LANGUAGE_ASSETS_PROVIDER = {
  provide: ALL_LANGUAGES_ASSETS_DATA,
  useValue: TESTING_ALL_LANGUAGE_ASSETS_DATA,
};
