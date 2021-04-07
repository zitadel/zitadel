import { Pipe, PipeTransform } from '@angular/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';

@Pipe({
    name: 'timestampToRetention',
})
export class TimestampToRetentionPipe implements PipeTransform {

    transform(value?: Duration.AsObject, ...args: unknown[]): unknown {
        if (value) {
            return this.retentionFromTimestamp(value);
        } else {
            return '';
        }
    }

    private retentionFromTimestamp(date: Duration.AsObject): any {
        if (date?.seconds !== undefined && date?.nanos !== undefined) {
            const ms = (date.seconds * 1000 + date.nanos / 1000 / 1000);
            const mins = ms / 1000 / 60;
            return mins / 60;
        }
    }
}

