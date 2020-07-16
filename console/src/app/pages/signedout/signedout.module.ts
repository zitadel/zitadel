import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { SharedModule } from 'src/app/modules/shared/shared.module';

import { SignedoutRoutingModule } from './signedout-routing.module';

@NgModule({
    declarations: [],
    imports: [
        CommonModule,
        SignedoutRoutingModule,
        SharedModule,
    ],
})
export class SignedoutModule { }
