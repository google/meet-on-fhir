import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot, UrlTree} from '@angular/router';
import {Observable} from 'rxjs';

import {LanguagesService} from '../languages.service';

@Injectable({providedIn: 'root'})
export class LanguageSelectorGuard implements CanActivate {
  constructor(
      private readonly languages: LanguagesService,
      private readonly router: Router) {}

  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot):
      Observable<boolean|UrlTree>|Promise<boolean|UrlTree>|boolean|UrlTree {
    if (this.languages.supportedLanguageCodes.length === 1) {
      this.router.navigateByUrl('/consent');
      return false;
    } else {
      return true;
    }
  }
}
