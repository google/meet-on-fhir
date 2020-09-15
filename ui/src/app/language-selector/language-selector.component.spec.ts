import {HarnessLoader} from '@angular/cdk/testing';
import {TestbedHarnessEnvironment} from '@angular/cdk/testing/testbed';
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {MatCardModule} from '@angular/material/card';
import {MatListModule} from '@angular/material/list';
import {MatListOptionHarness, MatSelectionListHarness} from '@angular/material/list/testing';
import {Router} from '@angular/router';
import {RouterTestingModule} from '@angular/router/testing';

import {getNativeName, LanguageCode} from '../i18n-helper';
import {LanguagesService} from '../languages.service';
import {TEST_ROUTES} from '../testing/routes';

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
    expect(selected).toBe(getNativeName(component.defaultLanguage));
  });

  it('should display all language options', async () => {
    const options = await loader.getAllHarnesses(MatListOptionHarness);
    const labels = options.map(option => option.getText());
    const labelTexts = (await Promise.all(labels)).map(l => l.trim());
    expect(labelTexts).toEqual(TEST_LANGUAGES.map(getNativeName));
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
