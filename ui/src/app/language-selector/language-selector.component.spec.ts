import {HarnessLoader} from '@angular/cdk/testing';
import {TestbedHarnessEnvironment} from '@angular/cdk/testing/testbed';
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {MatCardModule} from '@angular/material/card';
import {MatListModule} from '@angular/material/list';
import {MatListOptionHarness, MatSelectionListHarness} from '@angular/material/list/testing';
import {Router} from '@angular/router';
import {RouterTestingModule} from '@angular/router/testing';

import {getNativeName, LanguageCode} from '../i18n-helper';
import {ALL_LANGUAGE_ASSETS_PROVIDER, ALL_LANGUAGES_ASSETS_DATA} from '../i18n-strings';
import {LanguagesService} from '../languages.service';
import {TEST_ROUTES} from '../testing/routes';
import {TESTING_ALL_LANGUAGE_ASSETS_PROVIDER} from '../testing/testing-i18n-strings';

import {LanguageSelectorComponent} from './language-selector.component';

const TEST_LANGUAGES = [LanguageCode.Spanish, LanguageCode['English(US)']];

describe('LanguageSelectorComponent', () => {
  let component: LanguageSelectorComponent;
  let loader: HarnessLoader;
  let languages: LanguagesService;
  let fixture: ComponentFixture<LanguageSelectorComponent>;
  let router: Router;

  beforeEach(async(() => {
    TestBed
        .configureTestingModule({
          imports: [
            MatCardModule,
            MatListModule,
            RouterTestingModule.withRoutes(TEST_ROUTES),
          ],
          providers: [TESTING_ALL_LANGUAGE_ASSETS_PROVIDER],
          declarations: [LanguageSelectorComponent],
        })
        .compileComponents();
  }));

  beforeEach(() => {
    languages = TestBed.inject(LanguagesService);
    // make sure two options shows up in our test.
    spyOnProperty(languages, 'supportedLanguageCodes', 'get')
        .and.returnValue(TEST_LANGUAGES);

    fixture = TestBed.createComponent(LanguageSelectorComponent);
    loader = TestbedHarnessEnvironment.loader(fixture);
    component = fixture.componentInstance;
    router = TestBed.inject(Router);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should select default language', async () => {
    const selected = await (await loader.getHarness(MatListOptionHarness.with({
                       selected: true
                     }))).getText();
    expect(selected).toBe(getNativeName(languages.language));
  });

  it('should display all language options', async () => {
    const options = await loader.getAllHarnesses(MatListOptionHarness);
    const labels = options.map(option => option.getText());
    const labelTexts = (await Promise.all(labels)).map(l => l.trim());
    expect(labelTexts).toEqual(TEST_LANGUAGES.map(getNativeName));
  });


  it('should update strings after select a language', async () => {
    expect(
        fixture.debugElement.nativeElement.querySelector('button').textContent)
        .toBe('English');
    const select = await loader.getHarness(MatSelectionListHarness);
    await select.selectItems({text: getNativeName(LanguageCode.Spanish)});
    expect(
        fixture.debugElement.nativeElement.querySelector('button').textContent)
        .toBe('Spanish');
  });

  it('should submit selected language and navigate', async () => {
    spyOn(component, 'setLanguage').and.callThrough();
    spyOn(router, 'navigateByUrl').and.returnValue(Promise.resolve(true));

    const select = await loader.getHarness(MatSelectionListHarness);
    await select.selectItems({text: getNativeName(LanguageCode.Spanish)});

    fixture.debugElement.nativeElement.querySelector('button').click();

    expect(component.setLanguage).toHaveBeenCalledWith(LanguageCode.Spanish);
    expect(router.navigateByUrl).toHaveBeenCalledWith('/consent');
  });
});
