import { Component, Input } from '@angular/core';
import { OIDCAppType } from 'src/app/proto/generated/zitadel/app_pb';

@Component({
    selector: 'cnsl-app-card',
    templateUrl: './app-card.component.html',
    styleUrls: ['./app-card.component.scss'],
})
export class AppCardComponent {
    @Input() public outline: boolean = false;
    @Input() public type!: OIDCAppType;
    @Input() public isApiApp: boolean = false;
    public OIDCApplicationType: any = OIDCAppType;
}
