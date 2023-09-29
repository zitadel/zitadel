import { DataSource } from '@angular/cdk/collections';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { BehaviorSubject, from, Observable, of } from 'rxjs';
import { catchError, finalize, map } from 'rxjs/operators';
import { Membership } from 'src/app/proto/generated/zitadel/user_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';

export class MembershipsDataSource extends DataSource<Membership.AsObject> {
  public totalResult: number = 0;
  public viewTimestamp!: Timestamp.AsObject;
  public membershipsSubject: BehaviorSubject<Membership.AsObject[]> = new BehaviorSubject<Membership.AsObject[]>([]);
  private loadingSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading$: Observable<boolean> = this.loadingSubject.asObservable();

  constructor(
    private auth: GrpcAuthService,
    private service: ManagementService,
  ) {
    super();
  }

  public loadMemberships(userId: string, pageIndex: number, pageSize: number): void {
    const offset = pageIndex * pageSize;

    this.loadingSubject.next(true);

    from(this.service.listUserMemberships(userId, pageSize, offset))
      .pipe(
        map((resp) => {
          this.totalResult = resp.details?.totalResult || 0;
          if (resp.details?.viewTimestamp) {
            this.viewTimestamp = resp.details?.viewTimestamp;
          }
          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((members) => {
        this.membershipsSubject.next(members);
      });
  }

  public loadMyMemberships(pageIndex: number, pageSize: number): void {
    const offset = pageIndex * pageSize;

    this.loadingSubject.next(true);

    from(this.auth.listMyMemberships(pageSize, offset))
      .pipe(
        map((resp) => {
          this.totalResult = resp.details?.totalResult || 0;
          if (resp.details?.viewTimestamp) {
            this.viewTimestamp = resp.details?.viewTimestamp;
          }
          return resp.resultList;
        }),
        catchError(() => of([])),
        finalize(() => this.loadingSubject.next(false)),
      )
      .subscribe((members) => {
        this.membershipsSubject.next(members);
      });
  }

  public connect(): Observable<Membership.AsObject[]> {
    return this.membershipsSubject.asObservable();
  }

  public disconnect(): void {
    this.membershipsSubject.complete();
    this.loadingSubject.complete();
  }
}
