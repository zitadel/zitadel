import { Component, EventEmitter, Input, Output } from '@angular/core';
import { OIDCApplicationType } from 'src/app/proto/generated/management_pb';

export interface RadioItemAppType {
    type: OIDCApplicationType;
    titleI18nKey: string;
    descI18nKey: string;
    checked: boolean,
    disabled: boolean,
    prefix: string;
    background: string;
}

@Component({
    selector: 'app-type-radio',
    templateUrl: './app-type-radio.component.html',
    styleUrls: ['./app-type-radio.component.scss'],
})
export class AppTypeRadioComponent {
    selected: OIDCApplicationType = OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB;
    @Input() types!: RadioItemAppType[];
    @Output() selectedType: EventEmitter<OIDCApplicationType> = new EventEmitter();

    public emitChange(): void {
        this.selectedType.emit(this.selected);
    }
}