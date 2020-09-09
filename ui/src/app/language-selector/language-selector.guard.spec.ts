import {TestBed} from '@angular/core/testing';
import {ActivatedRouteSnapshot, Router, RouterStateSnapshot} from '@angular/router';
import {RouterTestingModule} from '@angular/router/testing';

import {LanguageCode} from '../i18n-helper';
import {LanguagesService} from '../languages.service';
import {TEST_ROUTES} from '../testing/routes';

import {LanguageSelectorGuard} from './language-selector.guard';

describe('LanguageSelectorGuard', () => {
  const next = {} as ActivatedRouteSnapshot;
  const state = {} as RouterStateSnapshot;
  let guard: LanguageSelectorGuard;
  let languages: LanguagesService;
  let router: Router;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule.withRoutes(TEST_ROUTES),
      ],
    });
    languages = TestBed.inject(LanguagesService);
    guard = TestBed.inject(LanguageSelectorGuard);
    router = TestBed.inject(Router);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });

  it('should activate if more than one supported languages', () => {
    spyOnProperty(languages, 'supportedLanguageCodes', 'get').and.returnValue([
      LanguageCode.Spanish, LanguageCode['English(US)']
    ]);
    expect(guard.canActivate(next, state)).toBe(true);
  });

  it('should navigate to consent page if only one supported language', () => {
    spyOn(router, 'navigateByUrl').and.returnValue(Promise.resolve(true));
    expect(guard.canActivate(next, state)).toBe(false);
    expect(router.navigateByUrl).toHaveBeenCalledWith('/consent');
  });
});
