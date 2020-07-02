import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { HttpLoaderFactory } from 'src/app/app.module';

import { UserGrantRoutingModule } from './user-grant-routing.module';
import { UserGrantComponent } from './user-grant.component';

@NgModule({
    declarations: [UserGrantComponent],
    imports: [
        CommonModule,
        UserGrantRoutingModule,
        MatIconModule,
        MatButtonModule,
        TranslateModule.forChild({
            loader: {
                provide: TranslateLoader,
                useFactory: HttpLoaderFactory,
                deps: [HttpClient],
            },
        }),
    ],
    schemas: [],
})
export class UserGrantModule { }
