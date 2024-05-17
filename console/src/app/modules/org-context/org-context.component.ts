import { SelectionModel } from '@angular/cdk/collections';
import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { BehaviorSubject, catchError, debounceTime, finalize, from, map, Observable, of, pipe, scan, take, tap } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgFieldName, OrgNameQuery, OrgQuery, OrgState, OrgStateQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ToastService } from 'src/app/services/toast.service';

const ORG_QUERY_LIMIT = 100;

@Component({
  selector: 'cnsl-org-context',
  templateUrl: './org-context.component.html',
  styleUrls: ['./org-context.component.scss'],
})
export class OrgContextComponent implements OnInit {
  public pinned: SelectionModel<Org.AsObject> = new SelectionModel<Org.AsObject>(true, []);

  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);

  public bottom: boolean = false;
  private _done: BehaviorSubject<any> = new BehaviorSubject(false);
  private _loading: BehaviorSubject<any> = new BehaviorSubject(false);
  public _orgs: BehaviorSubject<Org.AsObject[]> = new BehaviorSubject<Org.AsObject[]>([]);

  public orgs$: Observable<Org.AsObject[]> = this._orgs.pipe(
    map((orgs) => {
      return orgs.sort((left, right) => left.name.localeCompare(right.name));
    }),
    pipe(
      tap((orgs: Org.AsObject[]) => {
        this.pinned.clear();
        this.getPrefixedItem('pinned-orgs').then((stringifiedOrgs) => {
          if (stringifiedOrgs) {
            const orgIds: string[] = JSON.parse(stringifiedOrgs);
            const pinnedOrgs = orgs.filter((o) => orgIds.includes(o.id));
            pinnedOrgs.forEach((o) => this.pinned.select(o));
          }
        });
      }),
    ),
  );

  public filterControl: UntypedFormControl = new UntypedFormControl('');
  @Input() public org!: Org.AsObject;
  @ViewChild('input', { static: false }) input!: ElementRef;
  @Output() public closedCard: EventEmitter<void> = new EventEmitter();
  @Output() public setOrg: EventEmitter<Org.AsObject> = new EventEmitter();

  constructor(
    public authService: AuthenticationService,
    private auth: GrpcAuthService,
    private toast: ToastService,
  ) {
    this.filterControl.valueChanges.pipe(debounceTime(500)).subscribe((value) => {
      const filteredValues = this.loadOrgs(0, value.trim().toLowerCase());
      this.mapAndUpdate(filteredValues, true);
    });
  }

  public ngOnInit(): void {
    this.focusFilter();
    this.init();
  }

  public setActiveOrg(org: Org.AsObject) {
    this.setOrg.emit(org);
    this.closedCard.emit();
  }

  public onNearEndScroll(position: 'top' | 'bottom'): void {
    if (position === 'bottom') {
      this.more();
    }
  }

  public more(): void {
    const _cursor = this._orgs.getValue().length;
    let more: Promise<Org.AsObject[]> = this.loadOrgs(_cursor, '');
    this.mapAndUpdate(more);
  }

  public init(): void {
    let first: Promise<Org.AsObject[]> = this.loadOrgs(0);
    this.mapAndUpdate(first);
  }

  public loadOrgs(offset: number, filter?: string): Promise<Org.AsObject[]> {
    if (!filter) {
      const value = this.input?.nativeElement?.value;
      if (value) {
        filter = value;
      }
    }

    let query = new OrgQuery();
    const orgStateQuery = new OrgStateQuery();
    orgStateQuery.setState(OrgState.ORG_STATE_ACTIVE);
    query.setStateQuery(orgStateQuery);

    if (filter) {
      const orgNameQuery = new OrgNameQuery();
      orgNameQuery.setName(filter);
      orgNameQuery.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
      query.setNameQuery(orgNameQuery);
    }
    this.orgLoading$.next(true);
    return this.auth
      .listMyProjectOrgs(ORG_QUERY_LIMIT, offset, query ? [query] : undefined, OrgFieldName.ORG_FIELD_NAME_NAME, 'asc')
      .then((result) => {
        this.orgLoading$.next(false);
        return result.resultList;
      })
      .catch((error) => {
        this.orgLoading$.next(false);
        this.toast.showError(error);
        return [];
      });
  }

  private mapAndUpdate(col: Promise<Org.AsObject[]>, clear?: boolean): any {
    if (clear === false && (this._done.value || this._loading.value)) {
      return;
    }

    if (!this.bottom) {
      this._loading.next(true);

      return from(col)
        .pipe(
          take(1),
          tap((res: Org.AsObject[]) => {
            const current = this._orgs.getValue();
            if (clear) {
              this._orgs.next(res);
            } else {
              this._orgs.next([...current, ...res]);
            }

            this._loading.next(false);
            if (!res.length) {
              this._done.next(true);
            }
          }),
          catchError((_) => {
            this._loading.next(false);
            this.bottom = true;
            return of([]);
          }),
        )
        .subscribe();
    }
  }

  public closeCard(element: HTMLElement): void {
    if (!element.classList.contains('dontcloseonclick') && !element.classList.contains('mat-button-wrapper')) {
      this.closedCard.emit();
    }
  }

  private focusFilter(): void {
    setTimeout(() => {
      this.input.nativeElement.focus();
    }, 0);
  }

  public toggle(item: Org.AsObject, event: any): void {
    event.stopPropagation();
    this.pinned.toggle(item);

    this.setPrefixedItem('pinned-orgs', JSON.stringify(this.pinned.selected.map((o) => o.id)));
  }

  private async getPrefixedItem(key: string): Promise<string | null> {
    return localStorage.getItem(`${key}`);
  }

  private async setPrefixedItem(key: string, value: any): Promise<void> {
    return localStorage.setItem(`${key}`, value);
  }
}
