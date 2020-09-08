import {Component, OnInit} from '@angular/core';

import {SettingsService} from '../settings.service';

/**
 * This component allows user to selecte a preferred language.
 */
@Component({
  selector: 'app-language-selector',
  templateUrl: './language-selector.component.html',
  styleUrls: ['./language-selector.component.scss']
})
export class LanguageSelectorComponent implements OnInit {
  supportedLanguages: string[];

  defaultLanguage = 'English';

  constructor(private readonly settings: SettingsService) {
    this.supportedLanguages = this.settings.supportedLanguages;
  }

  ngOnInit(): void {}

  setLanguage(language: string): void {
    // TODO: navigate to the consent screen in selected language.
  }
}
