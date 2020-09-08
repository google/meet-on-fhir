import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree} from '@angular/router';
import {Observable} from 'rxjs';

import {SettingsService} from '../settings.service';

@Injectable({providedIn: 'root'})
export class LanguageSelectorGuard implements CanActivate {
  constructor(private readonly settings: SettingsService) {}

  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot):
      Observable<boolean|UrlTree>|Promise<boolean|UrlTree>|boolean|UrlTree {
    if (this.settings.supportedLanguages.length === 1) {
      // TODO: navigate to the consent screen.
      return false;
    } else {
      return true;
    }
  }
}
