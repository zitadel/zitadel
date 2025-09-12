import { Pipe, PipeTransform } from '@angular/core';
import { JavaScriptValue } from 'google-protobuf/google/protobuf/struct_pb';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';

@Pipe({
  name: 'topayload',
})
export class ToPayloadPipe implements PipeTransform {
  public transform(value: Event): JavaScriptValue | string {
    const pl = value.getPayload();
    if (pl) {
      return pl.toJavaScript();
    } else {
      return '';
    }
  }
}
