import {async, TestBed} from '@angular/core/testing';
import {RouterTestingModule} from '@angular/router/testing';

import {AppComponent} from './app.component';
import {SettingsService} from './settings.service';

describe('AppComponent', () => {
  let settings: SettingsService;

  beforeEach(async(() => {
    TestBed
        .configureTestingModule({
          imports: [RouterTestingModule],
          declarations: [AppComponent],
        })
        .compileComponents();

    settings = TestBed.inject(SettingsService);
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });

  it('should render toolbar', () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.nativeElement;
    expect(compiled.querySelector('.toolbar img').src)
        .toContain('assets/logo.png');
    expect(compiled.querySelector('.toolbar span').textContent)
        .toContain('Virtual waiting room');
    expect(compiled.querySelector('.org-info').textContent)
        .toContain(`${settings.organizationInfo.name}${
            settings.organizationInfo.phone}`);
  });

  it('should render footer', () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.nativeElement;
    expect(compiled.querySelector('.footer').textContent)
        .toContain(
            'Privacy' +
            'Terms' +
            'Help');
  });
});
