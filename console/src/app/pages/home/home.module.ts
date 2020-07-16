import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateModule } from '@ngx-translate/core';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';
import { SharedModule } from 'src/app/modules/shared/shared.module';

import { HomeRoutingModule } from './home-routing.module';
import { HomeComponent } from './home.component';



@NgModule({
    declarations: [HomeComponent],
    imports: [
        CommonModule,
        MatIconModule,
        HasRoleModule,
        HomeRoutingModule,
        MatButtonModule,
        TranslateModule,
        SharedModule,
    ],
})
export class HomeModule { }
