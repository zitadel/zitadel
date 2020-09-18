import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {MatPaginator, PageEvent} from "@angular/material/paginator";
import {MatTableDataSource} from "@angular/material/table";
import {
  ExternalIDPSearchResponse,
  ExternalIDPView,
} from "../../../../proto/generated/management_pb";
import {BehaviorSubject, Observable} from "rxjs";
import {ManagementService} from "../../../../services/mgmt.service";
import {ToastService} from "../../../../services/toast.service";
import {SelectionModel} from "@angular/cdk/collections";

@Component({
  selector: 'app-external-idps',
  templateUrl: './external-idps.component.html',
  styleUrls: ['./external-idps.component.scss']
})
export class ExternalIdpsComponent implements OnInit {
  @Input() userId!: string;
  @ViewChild(MatPaginator) public paginator!: MatPaginator;
  public externalIdpResult!: ExternalIDPSearchResponse.AsObject;
  public dataSource: MatTableDataSource<ExternalIDPView.AsObject> = new MatTableDataSource<ExternalIDPView.AsObject>();
  public selection: SelectionModel<ExternalIDPView.AsObject> = new SelectionModel<ExternalIDPView.AsObject>(true, []);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();
  @Input() public displayedColumns: string[] = [ 'idpConfigId', 'idpName', 'externalUserId', 'externalUserDisplayName'];

  constructor(private mgmtService: ManagementService,
              private toast: ToastService) { }

  ngOnInit(): void {
    this.getData(10, 0);
  }

  public isAllSelected(): boolean {
    const numSelected = this.selection.selected.length;
    const numRows = this.dataSource.data.length;
    return numSelected === numRows;
  }

  public masterToggle(): void {
    this.isAllSelected() ?
      this.selection.clear() :
      this.dataSource.data.forEach(row => this.selection.select(row));
  }

  public changePage(event: PageEvent): void {
    this.getData(event.pageSize, event.pageIndex * event.pageSize);
  }

  private async getData(limit: number, offset: number): Promise<void> {
    this.loadingSubject.next(true);

    this.mgmtService.SearchExternalIdps(this.userId, limit, offset).then(resp => {
      this.externalIdpResult = resp.toObject();
      this.dataSource.data = this.externalIdpResult.resultList;
      console.log(this.externalIdpResult.resultList);
      this.loadingSubject.next(false);
    }).catch((error: any) => {
      this.toast.showError(error);
      this.loadingSubject.next(false);
    });
  }

  public refreshPage(): void {
    this.getData(this.paginator.pageSize, this.paginator.pageIndex * this.paginator.pageSize);
  }
}
