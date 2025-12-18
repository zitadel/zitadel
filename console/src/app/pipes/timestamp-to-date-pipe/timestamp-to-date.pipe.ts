import { Pipe, PipeTransform } from "@angular/core";
import {
  Timestamp as BufTimestamp,
  TimestampSchema,
} from "@bufbuild/protobuf/wkt";
import { Timestamp } from "src/app/proto/generated/google/protobuf/timestamp_pb";
import { MessageInitShape } from "@bufbuild/protobuf";

@Pipe({
  name: "timestampToDate",
})
export class TimestampToDatePipe implements PipeTransform {
  transform(date: undefined): undefined;
  transform(date: BufTimestamp | Timestamp.AsObject): Date;
  transform(
    date: MessageInitShape<typeof TimestampSchema> | undefined
  ): Date | undefined;
  transform(
    date: BufTimestamp | Timestamp.AsObject | undefined
  ): Date | undefined;
  transform(
    date:
      | MessageInitShape<typeof TimestampSchema>
      | BufTimestamp
      | Timestamp.AsObject
      | undefined
  ): Date | undefined {
    if (date?.seconds !== undefined && date.nanos !== undefined) {
      return new Date(Number(date.seconds) * 1000 + date.nanos / 1000 / 1000);
    }
    return undefined;
  }
}
