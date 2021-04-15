import { Component, Input } from '@angular/core';

enum InfoSectionType {
    INFO = 'INFO',
    WARN = 'WARN',
}

@Component({
    selector: 'cnsl-info-section',
    templateUrl: './info-section.component.html',
    styleUrls: ['./info-section.component.scss'],
})
export class InfoSectionComponent {

    @Input() type: InfoSectionType = InfoSectionType.INFO;
}
