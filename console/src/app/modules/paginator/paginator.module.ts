import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatSelectModule } from "@angular/material/select";
import { TranslateModule } from "@ngx-translate/core";
import { LocalizedDatePipeModule } from "src/app/pipes/localized-date-pipe/localized-date-pipe.module";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";

import { FormFieldModule } from "../form-field/form-field.module";
import { PaginatorComponent } from "./paginator.component";

@NgModule({
  declarations: [PaginatorComponent],
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule,
    MatButtonModule,
    TimestampToDatePipe,
    FormFieldModule,
    MatSelectModule,
    LocalizedDatePipeModule,
  ],
  exports: [PaginatorComponent],
})
export class PaginatorModule {}
