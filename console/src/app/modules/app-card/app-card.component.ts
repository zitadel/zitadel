import { Component, Input } from '@angular/core';
import { OIDCApplicationType } from 'src/app/proto/generated/management_pb';

@Component({
    selector: 'cnsl-app-card',
    templateUrl: './app-card.component.html',
    styleUrls: ['./app-card.component.scss'],
})
export class AppCardComponent {
    @Input() public outline: boolean = false;
    @Input() public type!: OIDCApplicationType;
    @Input() public isApiApp: boolean = false;
    public OIDCApplicationType: any = OIDCApplicationType;
}
