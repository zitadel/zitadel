import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { TranslateModule } from '@ngx-translate/core';
import { InputModule } from 'src/app/modules/input/input.module';

import { DetailFormMachineComponent } from './detail-form-machine.component';


@NgModule({
    declarations: [
        DetailFormMachineComponent,
    ],
    imports: [
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        TranslateModule,
        InputModule,
        MatSelectModule,
        MatButtonModule,
        MatIconModule,
        TranslateModule,
    ],
    exports: [
        DetailFormMachineComponent,
    ],
})
export class DetailFormMachineModule { }
