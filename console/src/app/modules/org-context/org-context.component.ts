import { SelectionModel } from '@angular/cdk/collections';
import { Component, ElementRef, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { UntypedFormControl } from '@angular/forms';
import { BehaviorSubject, catchError, debounceTime, finalize, from, map, Observable, of, pipe, scan, take, tap } from 'rxjs';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Org, OrgFieldName, OrgNameQuery, OrgQuery, OrgState, OrgStateQuery } from 'src/app/proto/generated/zitadel/org_pb';
import { AuthenticationService } from 'src/app/services/authentication.service';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';

const ORG_QUERY_LIMIT = 4;

@Component({
  selector: 'cnsl-org-context',
  templateUrl: './org-context.component.html',
  styleUrls: ['./org-context.component.scss'],
})
export class OrgContextComponent implements OnInit {
  public pinned: SelectionModel<Org.AsObject> = new SelectionModel<Org.AsObject>(true, []);

  public orgLoading$: BehaviorSubject<any> = new BehaviorSubject(false);

  public bottom: boolean = false;
  private _cursor: BehaviorSubject<number> = new BehaviorSubject(0);
  private _done: BehaviorSubject<any> = new BehaviorSubject(false);
  private _loading: BehaviorSubject<any> = new BehaviorSubject(false);
  public _orgs: BehaviorSubject<Org.AsObject[]> = new BehaviorSubject<Org.AsObject[]>([]);

  public orgs$: Observable<Org.AsObject[]> = this._orgs.pipe(
    scan((acc, val) => {
      console.log(acc, val);
      return false ? val.concat(acc) : acc.concat(val);
    }),
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
  ) {
    this.filterControl.valueChanges.pipe(debounceTime(500)).subscribe((value) => {
      // this.orgs$.next([]);
      this.loadOrgs(0, value.trim().toLowerCase());
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
      console.log('more');
      this.more();
    }
  }

  public more(): void {
    const _cursor = this._cursor.getValue();
    const nextOffset = _cursor + ORG_QUERY_LIMIT;
    this._cursor.next(nextOffset);
    console.log('more', nextOffset);
    let more: Promise<Org.AsObject[]> = this.loadOrgs(nextOffset, '');
    this.mapAndUpdate(more);
  }

  public init(): void {
    let first: Promise<Org.AsObject[]>;
    first = this.loadOrgs(0);
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

    // this.orgLoading$.next(true);
    return this.auth
      .listMyProjectOrgs(ORG_QUERY_LIMIT, offset, query ? [query] : undefined, OrgFieldName.ORG_FIELD_NAME_NAME, 'asc')
      .then((result) => {
        return result.resultList;
      });
    // this.orgLoading$.next(false);
  }

  private mapAndUpdate(col: Promise<Org.AsObject[]>): any {
    if (this._done.value || this._loading.value) {
      return;
    }

    if (!this.bottom) {
      this._loading.next(true);

      return from(col)
        .pipe(
          take(1),
          tap((res: Org.AsObject[]) => {
            this._orgs.next(res);
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
