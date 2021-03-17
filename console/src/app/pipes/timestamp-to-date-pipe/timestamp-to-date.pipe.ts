import { Pipe, PipeTransform } from '@angular/core';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';

@Pipe({
    name: 'timestampToDate',
})
export class TimestampToDatePipe implements PipeTransform {

    transform(value: Timestamp.AsObject, ...args: unknown[]): unknown {
        return this.dateFromTimestamp(value);
    }

    private dateFromTimestamp(date: Timestamp.AsObject): any {
        if (date?.seconds && date?.nanos) {
            const ts: Date = new Date(date.seconds * 1000 + date.nanos / 1000 / 1000);
            return ts;
        }
    }
}

