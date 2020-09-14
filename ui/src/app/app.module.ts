import {NgModule} from '@angular/core';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {MatDividerModule} from '@angular/material/divider';
import {MatListModule} from '@angular/material/list';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {ConsentComponent} from './consent/consent.component';
import {ALL_LANGUAGE_ASSETS_PROVIDER} from './i18n-strings';
import {LanguageSelectorComponent} from './language-selector/language-selector.component';

@NgModule({
  declarations: [AppComponent, LanguageSelectorComponent, ConsentComponent],
  imports: [
    AppRoutingModule,
    BrowserModule,
    BrowserAnimationsModule,
    MatButtonModule,
    MatCardModule,
    MatDividerModule,
    MatListModule,
  ],
  providers: [
    ALL_LANGUAGE_ASSETS_PROVIDER,
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
