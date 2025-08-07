import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Observable, ReplaySubject } from 'rxjs';
import { Metadata as MetadataV2 } from '@zitadel/proto/zitadel/metadata_pb';
import { map, startWith } from 'rxjs/operators';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';

type StringMetadata = {
  key: string;
  value: string;
};

@Component({
  selector: 'cnsl-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.scss'],
})
export class MetadataComponent implements OnInit {
  @Input({ required: true }) public set metadata(metadata: (Metadata.AsObject | MetadataV2)[]) {
    this.metadata$.next(metadata);
  }
  @Input() public disabled: boolean = false;
  @Input() public loading: boolean = false;
  @Input({ required: true }) public description!: string;
  @Output() public editClicked: EventEmitter<void> = new EventEmitter();
  @Output() public refresh: EventEmitter<void> = new EventEmitter();

  public displayedColumns: string[] = ['key', 'value'];
  public metadata$ = new ReplaySubject<(Metadata.AsObject | MetadataV2)[]>(1);
  public dataSource$?: Observable<MatTableDataSource<StringMetadata>>;

  @ViewChild(MatSort) public sort!: MatSort;

  constructor() {}

  ngOnInit() {
    this.dataSource$ = this.metadata$.pipe(
      map((metadata) => {
        const decoder = new TextDecoder();
        return metadata.map(({ key, value }) => ({
          key,
          value: typeof value === 'string' ? value : decoder.decode(value),
        }));
      }),
      startWith([] as StringMetadata[]),
      map((metadata) => new MatTableDataSource(metadata)),
    );
  }
}
