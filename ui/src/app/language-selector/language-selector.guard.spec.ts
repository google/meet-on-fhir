import {TestBed} from '@angular/core/testing';
import {ActivatedRouteSnapshot, RouterStateSnapshot} from '@angular/router';

import {SettingsService} from '../settings.service';

import {LanguageSelectorGuard} from './language-selector.guard';

describe('LanguageSelectorGuard', () => {
  const next = {} as ActivatedRouteSnapshot;
  const state = {} as RouterStateSnapshot;
  let guard: LanguageSelectorGuard;
  let settings: SettingsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    settings = TestBed.inject(SettingsService);
    guard = TestBed.inject(LanguageSelectorGuard);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });

  it('should activate if more than one supported languages', () => {
    expect(guard.canActivate(next, state)).toBe(true);
  });

  it('should not activate if only one supported language', () => {
    spyOnProperty(settings, 'supportedLanguages', 'get').and.returnValue([
      'English'
    ]);
    expect(guard.canActivate(next, state)).toBe(false);
  });
});
