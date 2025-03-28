import { Pipe, PipeTransform } from '@angular/core';
import { Timestamp as BufTimestamp } from '@bufbuild/protobuf/wkt';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';

@Pipe({
  name: 'timestampToDate',
})
export class TimestampToDatePipe implements PipeTransform {
  transform(date: BufTimestamp | Timestamp.AsObject | undefined): Date | undefined {
    if (date?.seconds && date.nanos) {
      return new Date(Number(date.seconds) * 1000 + date.nanos / 1000 / 1000);
    }
    return undefined;
  }
}
