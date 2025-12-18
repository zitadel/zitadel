import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatTooltipModule } from "@angular/material/tooltip";
import { TranslateModule } from "@ngx-translate/core";
import { CopyToClipboardModule } from "src/app/directives/copy-to-clipboard/copy-to-clipboard.module";
import { LocalizedDatePipeModule } from "src/app/pipes/localized-date-pipe/localized-date-pipe.module";

import { MatDialogModule } from "@angular/material/dialog";
import { InfoSectionModule } from "../info-section/info-section.module";
import { ShowTokenDialogComponent } from "./show-token-dialog.component";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";

@NgModule({
  declarations: [ShowTokenDialogComponent],
  imports: [
    CommonModule,
    TranslateModule,
    MatDialogModule,
    InfoSectionModule,
    CopyToClipboardModule,
    MatButtonModule,
    MatTooltipModule,
    LocalizedDatePipeModule,
    TimestampToDatePipe,
  ],
})
export class ShowTokenDialogModule {}
