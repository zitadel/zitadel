import { Component, EventEmitter, Input, OnChanges, Output, SimpleChanges, ViewChild } from '@angular/core';
import { MatLegacyTable as MatTable, MatLegacyTableDataSource as MatTableDataSource } from '@angular/material/legacy-table';
import { MatSort } from '@angular/material/sort';
import { BehaviorSubject, Observable } from 'rxjs';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';

@Component({
  selector: 'cnsl-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.scss'],
})
export class MetadataComponent implements OnChanges {
  @Input() public metadata: Metadata.AsObject[] = [];
  @Input() public disabled: boolean = false;
  @Input() public loading: boolean = false;
  @Output() public editClicked: EventEmitter<void> = new EventEmitter();
  @Output() public refresh: EventEmitter<void> = new EventEmitter();

  public displayedColumns: string[] = ['key', 'value'];
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  @ViewChild(MatTable) public table!: MatTable<Metadata.AsObject>;
  @ViewChild(MatSort) public sort!: MatSort;
  public dataSource: MatTableDataSource<Metadata.AsObject> = new MatTableDataSource<Metadata.AsObject>([]);

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.metadata?.currentValue) {
      this.dataSource = new MatTableDataSource<Metadata.AsObject>(changes.metadata.currentValue);
    }
  }
}
