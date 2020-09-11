import {HarnessLoader} from '@angular/cdk/testing';
import {TestbedHarnessEnvironment} from '@angular/cdk/testing/testbed';
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {MatListModule} from '@angular/material/list';
import {MatListOptionHarness, MatSelectionListHarness} from '@angular/material/list/testing';

import {SettingsService} from '../settings.service';

import {LanguageSelectorComponent} from './language-selector.component';

describe('LanguageSelectorComponent', () => {
  let component: LanguageSelectorComponent;
  let loader: HarnessLoader;
  let settings: SettingsService;
  let fixture: ComponentFixture<LanguageSelectorComponent>;

  beforeEach(async(() => {
    TestBed
        .configureTestingModule({
          imports: [MatListModule],
          declarations: [LanguageSelectorComponent]
        })
        .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LanguageSelectorComponent);
    loader = TestbedHarnessEnvironment.loader(fixture);
    settings = TestBed.inject(SettingsService);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should select default language', async () => {
    const selected = await (await loader.getHarness(MatListOptionHarness.with({
                       selected: true
                     }))).getText();
    expect(selected).toBe(component.defaultLanguage);
  });

  it('should display all language options', async () => {
    const options = await loader.getAllHarnesses(MatListOptionHarness);
    const labels = options.map(option => option.getText());
    const labelTexts = (await Promise.all(labels)).map(l => l.trim());
    expect(labelTexts).toEqual(settings.supportedLanguages);
  });

  it('should submit selected language', async () => {
    spyOn(component, 'setLanguage');

    const select = await loader.getHarness(MatSelectionListHarness);
    await select.selectItems({text: 'Deutsch'});

    fixture.debugElement.nativeElement.querySelector('button').click();

    expect(component.setLanguage).toHaveBeenCalledWith('Deutsch');
  });
});
