import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';

import { AddIdpDialogModule } from './add-idp-dialog/add-idp-dialog.module';
import { LoginPolicyRoutingModule } from './login-policy-routing.module';
import { LoginPolicyComponent } from './login-policy.component';

@NgModule({
    declarations: [LoginPolicyComponent],
    imports: [
        LoginPolicyRoutingModule,
        CommonModule,
        FormsModule,
        MatInputModule,
        MatFormFieldModule,
        MatButtonModule,
        MatSlideToggleModule,
        MatIconModule,
        HasRolePipeModule,
        MatTooltipModule,
        TranslateModule,
        DetailLayoutModule,
        AddIdpDialogModule,
    ],
})
export class LoginPolicyModule { }
