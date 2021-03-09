import { Component, EventEmitter, Input, Output } from '@angular/core';
import {
    APIAuthMethodType,
    OIDCAuthMethodType,
    OIDCGrantType,
    OIDCResponseType,
} from 'src/app/proto/generated/zitadel/app_pb';

export interface RadioItemAuthType {
    key: string;
    titleI18nKey: string;
    descI18nKey: string;
    disabled: boolean,
    prefix: string;
    background: string;
    responseType?: OIDCResponseType;
    grantType?: OIDCGrantType;
    authMethod?: OIDCAuthMethodType;
    apiAuthMethod?: | APIAuthMethodType;
    recommended?: boolean;
    notRecommended?: boolean;
}

@Component({
    selector: 'app-auth-method-radio',
    templateUrl: './app-auth-method-radio.component.html',
    styleUrls: ['./app-auth-method-radio.component.scss'],
})
export class AppAuthMethodRadioComponent {
    @Input() current: string = '';
    @Input() selected: string = '';
    @Input() authMethods!: RadioItemAuthType[];
    @Input() isOIDC: boolean = false;
    @Output() selectedMethod: EventEmitter<string> = new EventEmitter();

    public emitChange(): void {
        this.selectedMethod.emit(this.selected);
    }
}