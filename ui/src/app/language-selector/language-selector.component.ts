import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';

import {getNativeName, LanguageCode} from '../i18n-helper';
import {LanguagesService} from '../languages.service';

/**
 * This component allows user to selecte a preferred language.
 */
@Component({
  selector: 'app-language-selector',
  templateUrl: './language-selector.component.html',
  styleUrls: ['./language-selector.component.scss']
})
export class LanguageSelectorComponent implements OnInit {
  supportedLanguages: [LanguageCode, string][];

  defaultLanguage = LanguageCode['English(US)'];

  constructor(
      readonly languages: LanguagesService, private readonly router: Router) {
    this.supportedLanguages =
        this.languages.supportedLanguageCodes.map(code => {
          return [code, getNativeName(code)];
        });
  }

  ngOnInit(): void {}

  setLanguage(language: LanguageCode): void {
    this.languages.selectedLanguage = language;
    this.router.navigateByUrl('/consent');
  }
}
