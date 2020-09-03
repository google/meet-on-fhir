import {Component} from '@angular/core';

import {OrganizationInfo, SettingsService} from './settings.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  organizationInfo: OrganizationInfo;

  constructor(private readonly settings: SettingsService) {
    this.organizationInfo = this.settings.organizationInfo;
  }
}
