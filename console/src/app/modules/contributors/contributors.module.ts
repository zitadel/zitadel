import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';

import { AvatarModule } from '../avatar/avatar.module';
import { ContributorsComponent } from './contributors.component';



@NgModule({
    declarations: [ContributorsComponent],
    imports: [
        CommonModule,
        AvatarModule,
        MatIconModule,
        MatTooltipModule,
        MatButtonModule,
    ],
    exports: [
        ContributorsComponent,
    ],
})
export class ContributorsModule { }
