import {Injectable} from '@angular/core';

import {LanguageCode, StringAssets} from './i18n-helper';
import {ALL_LANGUAGE_ASSETS} from './i18n-strings';

/**
 * A service which provides translated strings showing in the UI.
 * TODO: If this file size grows out of control, we should switch to use
 * Angular's standard i18n framework.
 */
@Injectable({providedIn: 'root'})
export class LanguagesService {
  selectedLanguage: LanguageCode = LanguageCode['English(US)'];

  private allLanguageCodes =
      // It is safe to use type cast here because the compiler would emit error
      // if ALL_LANGUAGE_ASSETS has a key which is not a LanguageCode.
      Object.keys(ALL_LANGUAGE_ASSETS).map(code => code as LanguageCode);

  get supportedLanguageCodes(): LanguageCode[] {
    return this.allLanguageCodes;
  }

  get stringAssets(): StringAssets {
    return ALL_LANGUAGE_ASSETS[this.selectedLanguage] ||
        ALL_LANGUAGE_ASSETS[LanguageCode['English(US)']];
  }
}
