import { BreakpointObserver } from '@angular/cdk/layout';
import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { BehaviorSubject, Observable, Subject, takeUntil } from 'rxjs';
import { map } from 'rxjs/operators';
import { StorageLocation, StorageService } from 'src/app/services/storage.service';

export enum MetaTab {
  DETAIL = 'DETAIL',
  ACTIVITY = 'ACTIVITY',
}

@Component({
  selector: 'cnsl-meta-layout',
  templateUrl: './meta-layout.component.html',
  styleUrls: ['./meta-layout.component.scss'],
})
export class MetaLayoutComponent implements OnInit, OnDestroy {
  @Input() changedTab$: BehaviorSubject<MetaTab> = new BehaviorSubject<MetaTab>(MetaTab.DETAIL);
  @Output() siteChanged: EventEmitter<MetaTab> = new EventEmitter();
  private destroy: Subject<void> = new Subject();

  constructor(private breakpointObserver: BreakpointObserver, private storageService: StorageService) {
    this.isSmallScreen$.subscribe((small) => (this.hidden = small));
  }
  public hidden: boolean = false;
  public isSmallScreen$: Observable<boolean> = this.breakpointObserver.observe('(max-width: 1000px)').pipe(
    map((result) => {
      return result.matches;
    }),
  );

  public ngOnInit(): void {
    const view = this.storageService.getItem<MetaTab>('MetaLayout', StorageLocation.local);
    if (view) {
      this.siteChanged.emit(view);
    }

    this.changedTab$.pipe(takeUntil(this.destroy)).subscribe((site) => {
      const view = this.storageService.getItem<MetaTab>('MetaLayout', StorageLocation.local);
      if (view !== site) {
        this.storageService.setItem('MetaLayout', site, StorageLocation.local);
      }
    });
  }

  public ngOnDestroy(): void {
    this.destroy.next();
    this.destroy.complete();
  }
}
