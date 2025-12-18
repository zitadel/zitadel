import { CommonModule } from "@angular/common";
import { NgModule } from "@angular/core";
import { LocalizedDatePipeModule } from "../localized-date-pipe/localized-date-pipe.module";
import { MilestonePipe } from "./milestonePipe";
import { TimestampToDatePipe } from "@/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe";

@NgModule({
  declarations: [MilestonePipe],
  imports: [CommonModule, TimestampToDatePipe, LocalizedDatePipeModule],
  exports: [MilestonePipe],
})
export class MilestonePipeModule {}
