import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, Observable } from 'rxjs';
import { GrantedProject } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

/**
 * Data source for the ProjectMembers view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class ProjectGrantsDataSource extends DataSource<GrantedProject.AsObject> {
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;
  public grantsSubject: BehaviorSubject<GrantedProject.AsObject[]> = new BehaviorSubject<GrantedProject.AsObject[]>([]);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  constructor(
    private mgmtService: ManagementService,
    private toast: ToastService,
  ) {
    super();
  }

  public loadGrants(projectId: string, pageIndex: number, pageSize: number, sortDirection?: string): void {
    const offset = pageIndex * pageSize;

    this.loadingSubject.next(true);
    this.mgmtService
      .listProjectGrants(projectId, pageSize, offset)
      .then((resp) => {
        if (resp.details?.totalResult) {
          this.totalResult = resp.details.totalResult;
        } else {
          this.totalResult = 0;
        }

        if (resp.details?.viewTimestamp) {
          this.viewTimestamp = resp.details?.viewTimestamp;
        }

        this.grantsSubject.next(resp.resultList);
        this.loadingSubject.next(false);
      })
      .catch((error) => {
        this.toast.showError(error);
        this.loadingSubject.next(false);
      });
  }

  /**
   * Connect this data source to the table. The table will only update when
   * the returned stream emits new items.
   * @returns A stream of the items to be rendered.
   */
  public connect(): Observable<GrantedProject.AsObject[]> {
    return this.grantsSubject.asObservable();
  }

  /**
   *  Called when the table is being destroyed. Use this function, to clean up
   * any open connections or free any held resources that were set up during connect.
   */
  public disconnect(): void {
    this.grantsSubject.complete();
    this.loadingSubject.complete();
  }
}
