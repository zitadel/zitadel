import { Component, EventEmitter, Input, Output } from '@angular/core';

import { OIDC, RadioItemIdpType } from '../idptypes';

@Component({
  selector: 'cnsl-idp-type-radio',
  templateUrl: './idp-type-radio.component.html',
  styleUrls: ['./idp-type-radio.component.scss'],
})
export class IdpTypeRadioComponent {
  @Input() selected: RadioItemIdpType = OIDC;
  @Input() types!: RadioItemIdpType[];
  @Output() selectedType: EventEmitter<RadioItemIdpType> = new EventEmitter();

  public emitChange(): void {
    this.selectedType.emit(this.selected);
  }
}
