import {TestBed} from '@angular/core/testing';
import {take} from 'rxjs/operators';

import {LanguageCode, StringAssets} from './i18n-helper';
import {LanguagesService} from './languages.service';
import {TESTING_ALL_LANGUAGE_ASSETS_DATA, TESTING_ALL_LANGUAGE_ASSETS_PROVIDER} from './testing/testing-i18n-strings';

describe('LanguagesService', () => {
  let service: LanguagesService;

  beforeEach(() => {
    TestBed.configureTestingModule(
        {providers: [TESTING_ALL_LANGUAGE_ASSETS_PROVIDER]});
    service = TestBed.inject(LanguagesService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should only support English by default', () => {
    expect(service.supportedLanguageCodes).toEqual([
      LanguageCode['English(US)'], LanguageCode.Spanish
    ]);
  });

  it('should return English assets if select a unknown language', async () => {
    service.language = LanguageCode.Amharic;
    const assets = await service.stringAssets.pipe(take(1)).toPromise();
    expect(assets).toEqual(
        (TESTING_ALL_LANGUAGE_ASSETS_DATA[LanguageCode['English(US)']]) as
        StringAssets);
  });
});
