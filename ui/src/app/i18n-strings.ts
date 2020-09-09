import {InjectionToken} from '@angular/core';

import {LanguageCode, StringAssets} from './i18n-helper';

// This type define makes sure the language code used in ALL_LANGUAGE_ASSETS is
// expected and no missing field in any translation.
type LanguageAssets = Partial<Record<LanguageCode, StringAssets>>;

/**
 * Strings used in UI in all supported languages.
 * We only provide English string. To add a new lanuage, do the following:
 * 1. copy and paste everything in 'en'.
 * 2. change all engilish strings to the correct translation.
 * 3. change the key to the correponding language code string. The string must
 * be one of the defined LanguageCode in i18-helper.ts.
 */
export const ALL_LANGUAGE_ASSETS: LanguageAssets = {
  // English
  'en': {
    submitButton: 'Submit',
    consentMessageTitle:
        'Before you enter the waiting room, please review the following information about televisits, via MyChart:',
    consentMessageBodyTexts: [
      'Please take a moment to see that you are in a private location where you cannot be unintentionally overheard. If you cannot be in a private location, you can decide at any time whether to continue this visit or to end it.',
      'You understand there are potential risks to this technology, including interruptions, unauthorized access and technical difficulties.  You or your provider may need to discontinue the televisit at any time.',
      'By continuing to participate in this telehealth visit, you are providing verbal consent for treatment.'
    ],
    languageSelect: 'Select your language:',
    consentButton: 'Continue to check-in'
  },
};
