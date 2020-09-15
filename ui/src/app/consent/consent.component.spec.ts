import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {MatCardModule} from '@angular/material/card';

import {ConsentComponent} from './consent.component';
import { TESTING_ALL_LANGUAGE_ASSETS_PROVIDER } from '../testing/testing-i18n-strings';

describe('ConsentComponent', () => {
  let component: ConsentComponent;
  let fixture: ComponentFixture<ConsentComponent>;

  beforeEach(async(() => {
    TestBed
        .configureTestingModule({
          imports: [MatCardModule],
          providers: [TESTING_ALL_LANGUAGE_ASSETS_PROVIDER],
          declarations: [ConsentComponent]
        })
        .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConsentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
