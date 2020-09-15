import {Inject, Injectable} from '@angular/core';
import {BehaviorSubject, Observable} from 'rxjs';
import {map} from 'rxjs/operators';

import {LanguageCode, StringAssets} from './i18n-helper';
import {ALL_LANGUAGES_ASSETS_DATA} from './i18n-strings';

/**
 * A service which provides translated strings showing in the UI.
 * TODO: If this file size grows out of control, we should switch to use
 * Angular's standard i18n framework.
 */
@Injectable({providedIn: 'root'})
export class LanguagesService {
  // We indirectly inject ALL_LANGUAGES_ASSETS for easier testing.
  constructor(@Inject(ALL_LANGUAGES_ASSETS_DATA)
              private readonly allLanguageAssets) {}

  private selectedLanguage =
      new BehaviorSubject<LanguageCode>(LanguageCode['English(US)']);

  private allLanguageCodes =
      // It is safe to use type cast here because the compiler would emit error
      // if ALL_LANGUAGE_ASSETS has a key which is not a LanguageCode.
      Object.keys(this.allLanguageAssets).map(code => code as LanguageCode);

  set language(language: LanguageCode) {
    this.selectedLanguage.next(language);
  }

  get language(): LanguageCode {
    return this.selectedLanguage.value;
  }

  get supportedLanguageCodes(): LanguageCode[] {
    return this.allLanguageCodes;
  }

  get stringAssets(): Observable<StringAssets> {
    return this.selectedLanguage.pipe(map(code => {
      return this.allLanguageAssets[code] ||
          this.allLanguageAssets[LanguageCode['English(US)']];
    }));
  }
}
