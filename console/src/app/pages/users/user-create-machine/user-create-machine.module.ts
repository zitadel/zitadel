import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';
import { InputModule } from 'src/app/modules/input/input.module';

import { UserCreateMachineRoutingModule } from './user-create-machine-routing.module';
import { UserCreateMachineComponent } from './user-create-machine.component';



@NgModule({
    declarations: [UserCreateMachineComponent],
    imports: [
        UserCreateMachineRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatSelectModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatProgressBarModule,
        MatCheckboxModule,
        MatTooltipModule,
        TranslateModule,
        DetailLayoutModule,
        InputModule,
    ],
})
export class UserCreateMachineModule { }
