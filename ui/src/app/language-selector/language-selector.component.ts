import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';

import {getNativeName, LanguageCode} from '../i18n-helper';
import {LanguagesService} from '../languages.service';

interface LanguageOption {
  code: LanguageCode;
  name: string;
}

/**
 * This component allows user to selecte a preferred language.
 */
@Component({
  selector: 'app-language-selector',
  templateUrl: './language-selector.component.html',
  styleUrls: ['./language-selector.component.scss']
})
export class LanguageSelectorComponent implements OnInit {
  supportedLanguages: LanguageOption[];

  defaultLanguage = LanguageCode['English(US)'];

  constructor(
      readonly languages: LanguagesService, private readonly router: Router) {
    this.supportedLanguages =
        this.languages.supportedLanguageCodes.map(code => {
          return {code, name: getNativeName(code)};
        });
  }

  ngOnInit(): void {}

  setLanguage(language: LanguageCode): void {
    this.languages.selectedLanguage = language;
    this.router.navigateByUrl('/consent');
  }
}
