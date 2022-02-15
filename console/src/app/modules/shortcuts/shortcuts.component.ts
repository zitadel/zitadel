import { CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { Component, OnDestroy } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { StorageLocation, StorageService } from 'src/app/services/storage.service';

import {
  LOGIN_POLICY,
  LOGIN_TEXTS_POLICY,
  MESSAGE_TEXTS_POLICY,
  POLICIES,
  PRIVATELABEL_POLICY,
} from '../policy-grid/policies';

export interface ShortcutItem {
  title?: string;
  desc?: string;
  i18nTitle?: string;
  i18nDesc?: string;
  routerLink: any;
  withRole: string[];
  icon?: string;
  svgIcon?: string;
  avatarSrc?: string;
  color?: string;
  disabled?: boolean;
}

const PROFILE_SHORTCUT: ShortcutItem = {
  routerLink: ['/users', 'me'],
  i18nTitle: 'USER.TITLE',
  icon: 'las la-cog',
  withRole: [''],
  disabled: true,
};

const CREATE_ORG: ShortcutItem = {
  i18nTitle: 'ORG.PAGES.CREATE',
  routerLink: ['/org', 'create'],
  withRole: ['org.create', 'iam.write'],
  icon: 'las la-plus',
  disabled: true,
};

const CREATE_PROJECT: ShortcutItem = {
  i18nTitle: 'PROJECT.PAGES.CREATE',
  routerLink: ['/projects', 'create'],
  withRole: ['project.create'],
  icon: 'las la-plus',
  disabled: true,
};

@Component({
  selector: 'cnsl-shortcuts',
  templateUrl: './shortcuts.component.html',
  styleUrls: ['./shortcuts.component.scss'],
})
export class ShortcutsComponent implements OnDestroy {
  public main: ShortcutItem[] = [PROFILE_SHORTCUT, CREATE_ORG];
  public secondary: ShortcutItem[] = [];
  public third: ShortcutItem[] = [];

  public ALL_POLICIES = POLICIES.map((p) => {
    const policy: ShortcutItem = {
      i18nTitle: p.i18nTitle,
      i18nDesc: p.i18nDesc,
      routerLink: p.orgRouterLink,
      withRole: p.orgWithRole,
      icon: p.icon ?? '',
      svgIcon: p.svgIcon ?? '',
      color: p.color ?? '',
      disabled: false,
    };
    return policy;
  });

  public all: ShortcutItem[] = this.ALL_POLICIES;

  private destroy$: Subject<void> = new Subject();
  public editState: boolean = false;
  constructor(private storageService: StorageService, private auth: GrpcAuthService) {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      this.loadShortcuts(org);
    }

    this.auth.activeOrgChanged.pipe(takeUntil(this.destroy$)).subscribe((org) => {
      this.loadShortcuts(org);
    });
  }

  public loadShortcuts(org: Org.AsObject): void {
    ['main', 'secondary', 'third'].map((listName) => {
      const shortcuts = this.storageService.getItem(`shortcuts:${listName}:${org.id}`, StorageLocation.local);
      if (shortcuts) {
        const parsed = JSON.parse(shortcuts);
        if (parsed) {
          switch (listName) {
            case 'main':
              const f = parsed.filter(
                (shortcut: ShortcutItem) =>
                  JSON.stringify(shortcut) !== JSON.stringify(PROFILE_SHORTCUT) &&
                  JSON.stringify(shortcut) !== JSON.stringify(CREATE_ORG),
              );
              this.main = [PROFILE_SHORTCUT, CREATE_ORG, ...f];
              break;
            case 'secondary':
              this.secondary = parsed;
              break;
            case 'third':
              this.third = parsed;
              break;
          }
          this.organizeAllItems();
        }
      } else {
        switch (listName) {
          case 'main':
            this.main = [PROFILE_SHORTCUT, CREATE_ORG, CREATE_PROJECT];
            break;
          case 'secondary':
            this.secondary = [LOGIN_POLICY, PRIVATELABEL_POLICY].map((p) => {
              const policy: ShortcutItem = {
                i18nTitle: p.i18nTitle,
                i18nDesc: p.i18nDesc,
                routerLink: p.orgRouterLink,
                withRole: p.orgWithRole,
                icon: p.icon ?? '',
                color: p.color ?? '',
                disabled: false,
              };
              return policy;
            });
            break;
          case 'third':
            this.third = [LOGIN_TEXTS_POLICY, MESSAGE_TEXTS_POLICY].map((p) => {
              const policy: ShortcutItem = {
                i18nTitle: p.i18nTitle,
                i18nDesc: p.i18nDesc,
                routerLink: p.orgRouterLink,
                withRole: p.orgWithRole,
                icon: p.icon ?? '',
                color: p.color ?? '',
                disabled: false,
              };
              return policy;
            });
            break;
        }
        this.organizeAllItems();
      }
    });
  }

  private organizeAllItems(): void {
    const list = [this.main, this.secondary, this.third].flat();
    const filtered = this.ALL_POLICIES.filter((p) => !list.find((l) => l.i18nTitle === p.i18nTitle));
    this.all = filtered;
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public drop(event: CdkDragDrop<ShortcutItem[]>, listName: string) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
      this.saveStateToStorage();
    } else {
      transferArrayItem(event.previousContainer.data, event.container.data, event.previousIndex, event.currentIndex);
      this.saveStateToStorage();
    }
  }

  public saveStateToStorage(): void {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      this.storageService.setItem(`shortcuts:main:${org.id}`, JSON.stringify(this.main), StorageLocation.local);
      this.storageService.setItem(`shortcuts:secondary:${org.id}`, JSON.stringify(this.secondary), StorageLocation.local);
      this.storageService.setItem(`shortcuts:third:${org.id}`, JSON.stringify(this.third), StorageLocation.local);
    }
  }

  public reset(): void {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      ['main', 'secondary', 'third'].map((listName) => {
        this.storageService.removeItem(`shortcuts:${listName}:${org.id}`, StorageLocation.local);
      });

      this.loadShortcuts(org);
    }
  }
}
