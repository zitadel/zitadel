import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatRippleModule } from '@angular/material/core';
import { TranslateModule } from '@ngx-translate/core';

import { AppAuthMethodRadioComponent } from './app-auth-method-radio/app-auth-method-radio.component';
import { AppTypeRadioComponent } from './app-type-radio/app-type-radio.component';

@NgModule({
    declarations: [
        AppTypeRadioComponent,
        AppAuthMethodRadioComponent,
    ],
    imports: [
        CommonModule,
        FormsModule,
        MatRippleModule,
        TranslateModule,
    ],
    exports: [
        AppAuthMethodRadioComponent,
        AppTypeRadioComponent,
    ],
})
export class AppRadioModule { }

