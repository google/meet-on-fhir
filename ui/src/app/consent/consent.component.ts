import {Component} from '@angular/core';
import {LanguagesService} from '../languages.service';

@Component({
  selector: 'app-consent',
  templateUrl: './consent.component.html',
  styleUrls: ['./consent.component.scss']
})
export class ConsentComponent {
  constructor(readonly languages: LanguagesService) {}
}
