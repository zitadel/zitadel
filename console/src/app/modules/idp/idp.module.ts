import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatChipsModule } from '@angular/material/chips';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';

import { IdpRoutingModule } from './idp-routing.module';
import { IdpComponent } from './idp.component';

@NgModule({
    declarations: [IdpComponent],
    imports: [
        CommonModule,
        IdpRoutingModule,
        FormsModule,
        ReactiveFormsModule,
        InputModule,
        MatButtonModule,
        MatIconModule,
        MatTooltipModule,
        MatSelectModule,
        TranslateModule,
        MatCheckboxModule,
        MatChipsModule,
        DetailLayoutModule,
    ],
})
export class IdpModule { }
