import { Component } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { MetaTab } from '../meta-layout.component';

@Component({
  selector: 'cnsl-remembered-tab',
  templateUrl: './remembered-tab.component.html',
  styleUrls: ['./remembered-tab.component.scss'],
})
export class RememberedTabComponent {
  public MetaTab: any = MetaTab;
  public metaTab!: MetaTab;
  public selectedMetaTab$: BehaviorSubject<MetaTab> = new BehaviorSubject<MetaTab>(MetaTab.DETAIL);

  constructor() {}

  public setTab(site: MetaTab): void {
    this.metaTab = site;
    this.selectedMetaTab$.next(site);
  }

  public selectTab(index: number): void {
    this.selectedMetaTab$.next(index === 0 ? MetaTab.DETAIL : MetaTab.ACTIVITY);
  }
}
