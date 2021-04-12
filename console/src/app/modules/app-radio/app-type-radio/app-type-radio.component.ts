import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RadioItemAppType, WEB_TYPE } from 'src/app/pages/projects/apps/authtypes';

@Component({
    selector: 'app-type-radio',
    templateUrl: './app-type-radio.component.html',
    styleUrls: ['./app-type-radio.component.scss'],
})
export class AppTypeRadioComponent {
    @Input() selected: RadioItemAppType = WEB_TYPE;
    @Input() types!: RadioItemAppType[];
    @Output() selectedType: EventEmitter<RadioItemAppType> = new EventEmitter();

    public emitChange(): void {
        this.selectedType.emit(this.selected);
    }
}
