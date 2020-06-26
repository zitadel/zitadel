import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';

import { AvatarComponent } from './avatar.component';



@NgModule({
    declarations: [AvatarComponent],
    imports: [
        CommonModule,
    ],
    exports: [
        AvatarComponent,
    ],
})
export class AvatarModule { }

