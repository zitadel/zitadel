import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';

@Component({
  selector: 'cnsl-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.scss'],
})
export class MetadataComponent {
  @Input() public metadata: Metadata.AsObject[] = [];
  @Input() public disabled: boolean = false;
  @Input() public loading: boolean = false;
  @Output() public editClicked: EventEmitter<void> = new EventEmitter();
  @Output() public refresh: EventEmitter<void> = new EventEmitter();

  constructor() {}
}
