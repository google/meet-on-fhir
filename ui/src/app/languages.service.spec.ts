import {TestBed} from '@angular/core/testing';

import {LanguageCode} from './i18n-helper';
import {ALL_LANGUAGE_ASSETS} from './i18n-strings';
import {LanguagesService} from './languages.service';

describe('LanguagesService', () => {
  let service: LanguagesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LanguagesService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should only support English by default', () => {
    expect(service.supportedLanguageCodes).toEqual([
      LanguageCode['English(US)']
    ]);
  });

  it('should return English assets if select a unknown language', () => {
    service.selectedLanguage = LanguageCode.Amharic;
    expect(service.stringAssets)
        .toEqual(ALL_LANGUAGE_ASSETS[LanguageCode['English(US)']]);
  });
});
