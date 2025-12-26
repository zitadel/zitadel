import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { TranslateModule } from "@ngx-translate/core";
import { LocalizedDatePipeModule } from "src/app/pipes/localized-date-pipe/localized-date-pipe.module";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";
import { InfoSectionModule } from "../info-section/info-section.module";

import { MatDialogModule } from "@angular/material/dialog";
import { ShowKeyDialogComponent } from "./show-key-dialog.component";

@NgModule({
  declarations: [ShowKeyDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    MatButtonModule,
    LocalizedDatePipeModule,
    MatDialogModule,
    InfoSectionModule,
    TimestampToDatePipe,
  ],
})
export class ShowKeyDialogModule {}
