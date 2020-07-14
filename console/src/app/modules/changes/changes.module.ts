import { ScrollingModule } from '@angular/cdk/scrolling';
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { TranslateModule } from '@ngx-translate/core';
import { ScrollableModule } from 'src/app/directives/scrollable/scrollable.module';
import { HasRolePipeModule } from 'src/app/pipes/has-role-pipe.module';
import { PipesModule } from 'src/app/pipes/pipes.module';

import { ChangesComponent } from './changes.component';



@NgModule({
    declarations: [
        ChangesComponent,
    ],
    imports: [
        CommonModule,
        ScrollableModule,
        MatProgressSpinnerModule,
        TranslateModule,
        PipesModule,
        HasRolePipeModule,
        ScrollingModule,
    ],
    exports: [
        ChangesComponent,
        ScrollableModule,
    ],
})
export class ChangesModule { }
