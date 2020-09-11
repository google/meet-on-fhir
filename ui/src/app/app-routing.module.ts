import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {LanguageSelectorComponent} from './language-selector/language-selector.component';
import {LanguageSelectorGuard} from './language-selector/language-selector.guard';

const routes: Routes = [{
  path: 'select-language',
  component: LanguageSelectorComponent,
  canActivate: [LanguageSelectorGuard]
}];

@NgModule({imports: [RouterModule.forRoot(routes)], exports: [RouterModule]})
export class AppRoutingModule {
}
