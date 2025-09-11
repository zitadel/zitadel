import { Pipe, PipeTransform } from '@angular/core';
import { Struct } from 'google-protobuf/google/protobuf/struct_pb';
import { Event } from 'src/app/proto/generated/zitadel/event_pb';

@Pipe({
  name: 'toobject',
})
export class ToObjectPipe implements PipeTransform {
  public transform(value: Event | Struct): any {
    return value.toObject();
  }
}
