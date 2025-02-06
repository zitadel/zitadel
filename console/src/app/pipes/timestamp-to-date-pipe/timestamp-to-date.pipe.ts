import { Pipe, PipeTransform } from '@angular/core';
import { Timestamp as ConnectTimestamp } from '@bufbuild/protobuf/wkt';
import { Timestamp } from 'src/app/proto/generated/google/protobuf/timestamp_pb';

@Pipe({
  name: 'timestampToDate',
})
export class TimestampToDatePipe implements PipeTransform {
  transform(value: ConnectTimestamp | Timestamp.AsObject, ...args: unknown[]): unknown {
    return this.dateFromTimestamp(value);
  }

  private dateFromTimestamp(date: ConnectTimestamp | Timestamp.AsObject): any {
    if (date?.seconds !== undefined && date?.nanos !== undefined) {
      return new Date(Number(date.seconds) * 1000 + date.nanos / 1000 / 1000);
    }
  }
}
