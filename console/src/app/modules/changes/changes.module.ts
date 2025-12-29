import { ScrollingModule } from "@angular/cdk/scrolling";
import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatTooltipModule } from "@angular/material/tooltip";
import { TranslateModule } from "@ngx-translate/core";
import { HasRolePipeModule } from "src/app/pipes/has-role-pipe/has-role-pipe.module";
import { LocalizedDatePipeModule } from "src/app/pipes/localized-date-pipe/localized-date-pipe.module";
import { AvatarModule } from "../avatar/avatar.module";

import { ChangesComponent } from "./changes.component";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";

@NgModule({
  declarations: [ChangesComponent],
  imports: [
    CommonModule,
    MatProgressSpinnerModule,
    TranslateModule,
    MatIconModule,
    MatButtonModule,
    HasRolePipeModule,
    ScrollingModule,
    LocalizedDatePipeModule,
    TimestampToDatePipe,
    MatTooltipModule,
    AvatarModule,
  ],
  exports: [ChangesComponent],
})
export class ChangesModule {}
