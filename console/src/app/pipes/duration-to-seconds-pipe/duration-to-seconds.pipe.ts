import { Pipe, PipeTransform } from '@angular/core';
import { Duration } from 'google-protobuf/google/protobuf/duration_pb';

@Pipe({
  name: 'durationToSeconds',
})
export class DurationToSecondsPipe implements PipeTransform {
  transform(value?: Duration.AsObject, ...args: unknown[]): unknown {
    if (value) {
      return this.durationToSeconds(value);
    } else {
      return '';
    }
  }

  private durationToSeconds(date: Duration.AsObject): any {
    if (date?.seconds !== undefined && date?.nanos !== undefined) {
      const ms = date.seconds * 1000 + date.nanos / 1000 / 1000;
      const secs = ms / 1000;
      return `${secs.toFixed(2)} sec`;
    }
  }
}
