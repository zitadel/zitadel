import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSelectModule } from '@angular/material/select';
import { MatTooltipModule } from '@angular/material/tooltip';
import { TranslateModule } from '@ngx-translate/core';
import { DetailLayoutModule } from 'src/app/modules/detail-layout/detail-layout.module';

import { UserCreateMachineRoutingModule } from './user-create-machine-routing.module';
import { UserCreateMachineComponent } from './user-create-machine.component';



@NgModule({
    declarations: [UserCreateMachineComponent],
    imports: [
        UserCreateMachineRoutingModule,
        CommonModule,
        FormsModule,
        ReactiveFormsModule,
        MatInputModule,
        MatFormFieldModule,
        MatSelectModule,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
        MatProgressBarModule,
        MatCheckboxModule,
        MatTooltipModule,
        TranslateModule,
        DetailLayoutModule,
    ],
})
export class UserCreateMachineModule { }
