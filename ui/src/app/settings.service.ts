import {Injectable} from '@angular/core';

/**
 * The organization (e.g clinic, hospital) information which displayed in the
 * waiting room ui.
 */
export interface OrganizationInfo {
  name: string;
  phone: string;
}

/**
 * The organization (e.g clinic, hospital) information which displayed in the
 * waiting room ui.
 */
@Injectable({providedIn: 'root'})
export class SettingsService {
  constructor() {}

  get organizationInfo(): OrganizationInfo {
    // TOOD: fetch these settings from the backend.
    return {
      name: 'Your Organization Name',
      phone: '(123)123-1234',
    } as OrganizationInfo;
  }
}
