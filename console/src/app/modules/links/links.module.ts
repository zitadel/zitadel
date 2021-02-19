import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LinksComponent } from './links.component';
import { TranslateModule } from '@ngx-translate/core';
import { RouterModule } from '@angular/router';
import { MatButton, MatButtonModule } from '@angular/material/button';
import { HasRoleModule } from 'src/app/directives/has-role/has-role.module';



@NgModule({
    declarations: [LinksComponent],
    imports: [
        CommonModule,
        TranslateModule,
        RouterModule,
        MatButtonModule,
        HasRoleModule,
    ],
    exports: [
        LinksComponent,
    ]
})
export class LinksModule { }
