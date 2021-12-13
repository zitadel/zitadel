import { CdkDragDrop, CdkDropList, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { Component, OnDestroy } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { StorageLocation, StorageService } from 'src/app/services/storage.service';

import { GridPolicy, POLICIES } from '../policy-grid/policies';

@Component({
  selector: 'cnsl-shortcuts',
  templateUrl: './shortcuts.component.html',
  styleUrls: ['./shortcuts.component.scss'],
})
export class ShortcutsComponent implements OnDestroy {
  public shortcuts: GridPolicy[] = [];
  public POLICIES: GridPolicy[] = POLICIES;

  public all = POLICIES; // ['Get up', 'Brush teeth', 'Take a shower', 'Check e-mail', 'Walk dog'];

  private destroy$: Subject<void> = new Subject();
  public editState: boolean = false;
  constructor(private storageService: StorageService, private auth: GrpcAuthService) {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      this.loadShortcuts(org);
    }

    this.auth.activeOrgChanged.pipe(takeUntil(this.destroy$)).subscribe((org) => {
      console.log(org.name, org.id);
      this.loadShortcuts(org);
    });
  }

  public loadShortcuts(org: Org.AsObject): void {
    const shortcuts = this.storageService.getItem(`shortcuts:${org.id}`, StorageLocation.local);
    if (shortcuts) {
      const parsed = JSON.parse(shortcuts);
      if (parsed) {
        this.organizeItems(parsed);
      }
    } else {
      this.organizeItems([]);
    }
  }

  private organizeItems(list: GridPolicy[]): void {
    console.log(list);
    this.shortcuts = list;
    const filtered = POLICIES.filter((p) => !list.find((l) => l.i18nTitle === p.i18nTitle));
    this.all = filtered;
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public drop(event: CdkDragDrop<GridPolicy[]>) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
      this.saveToStorage(event.container);
    } else {
      transferArrayItem(event.previousContainer.data, event.container.data, event.previousIndex, event.currentIndex);
      this.saveToStorage(event.container);
    }
  }

  public saveToStorage(list: CdkDropList): void {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      this.storageService.setItem(`shortcuts:${org.id}`, JSON.stringify(list.data), StorageLocation.local);
    }
  }
}
