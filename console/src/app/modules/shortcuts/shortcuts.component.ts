import { CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { Component, effect, OnDestroy } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageLocation, StorageService } from 'src/app/services/storage.service';

import { NewOrganizationService } from '../../services/new-organization.service';

export interface SettingLinks {
  i18nTitle: string;
  i18nDesc: string;
  iamRouterLink: any;
  orgRouterLink?: any;
  queryParams: any;
  iamWithRole?: string[];
  orgWithRole?: string[];
  icon?: string;
  svgIcon?: string;
  color: string;
}

export const LOGIN_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.LOGIN',
  i18nDesc: 'POLICY.LOGIN_POLICY.DESCRIPTION',
  iamRouterLink: ['/settings'],
  orgRouterLink: ['/org-settings'],
  queryParams: { id: 'login' },
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-sign-in-alt',
  color: 'green',
};

export const APPEARANCE_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.APPEARANCE',
  i18nDesc: 'POLICY.PRIVATELABELING.DESCRIPTION',
  iamRouterLink: ['/settings'],
  orgRouterLink: ['/org-settings'],
  queryParams: { id: 'branding' },
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-swatchbook',
  color: 'blue',
};

export const PRIVACY_POLICY: SettingLinks = {
  i18nTitle: 'DESCRIPTIONS.SETTINGS.PRIVACY_POLICY.TITLE',
  i18nDesc: 'POLICY.PRIVACY_POLICY.DESCRIPTION',
  iamRouterLink: ['/settings'],
  orgRouterLink: ['/org-settings'],
  queryParams: { id: 'privacypolicy' },
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-file-contract',
  color: 'black',
};

export const NOTIFICATION_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.NOTIFICATIONS',
  i18nDesc: 'SETTINGS.LIST.NOTIFICATIONS_DESC',
  iamRouterLink: ['/settings'],
  queryParams: { id: 'smtpprovider' },
  iamWithRole: ['iam.policy.read'],
  icon: 'las la-bell',
  color: 'red',
};

export const SETTINGLINKS: SettingLinks[] = [LOGIN_GROUP, APPEARANCE_GROUP, PRIVACY_POLICY, NOTIFICATION_GROUP];

export interface ShortcutItem {
  id: string;
  type: ShortcutType;
  title?: string;
  desc?: string;
  i18nTitle?: string;
  i18nDesc?: string;
  routerLink: any;
  queryParams?: any;
  withRole?: string[];
  icon?: string;
  label?: string;
  svgIcon?: string;
  avatarSrc?: string;
  color?: string;
  disabled?: boolean;
  state?: ProjectState;
}

export enum ShortcutType {
  ROUTE,
  POLICY,
  PROJECT,
}

const PROFILE_SHORTCUT: ShortcutItem = {
  id: 'profile',
  type: ShortcutType.ROUTE,
  routerLink: ['/users', 'me'],
  i18nTitle: 'USER.TITLE',
  icon: 'las la-cog',
  withRole: [''],
  disabled: false,
};

const CREATE_ORG: ShortcutItem = {
  id: 'create_org',
  type: ShortcutType.ROUTE,
  i18nTitle: 'ORG.PAGES.CREATE',
  routerLink: ['/orgs', 'create'],
  withRole: ['org.create', 'iam.write'],
  icon: 'las la-plus',
  disabled: false,
};

const CREATE_PROJECT: ShortcutItem = {
  id: 'create_project',
  type: ShortcutType.ROUTE,
  i18nTitle: 'PROJECT.PAGES.CREATE',
  routerLink: ['/projects', 'create'],
  withRole: ['project.create'],
  icon: 'las la-plus',
  disabled: false,
};

const CREATE_USER: ShortcutItem = {
  id: 'create_user',
  type: ShortcutType.ROUTE,
  i18nTitle: 'USER.CREATE.TITLE',
  routerLink: ['/users', 'create'],
  withRole: ['user.write'],
  icon: 'las la-plus',
  disabled: false,
};

@Component({
  selector: 'cnsl-shortcuts',
  templateUrl: './shortcuts.component.html',
  styleUrls: ['./shortcuts.component.scss'],
  standalone: false,
})
export class ShortcutsComponent implements OnDestroy {
  public orgId!: string;

  public main: ShortcutItem[] = [];
  public secondary: ShortcutItem[] = [];
  public third: ShortcutItem[] = [];

  public ALL_SHORTCUTS: ShortcutItem[] = [];
  public all: ShortcutItem[] = [];

  private destroy$: Subject<void> = new Subject();
  public editState: boolean = false;
  public ProjectState: any = ProjectState;

  constructor(
    private storageService: StorageService,
    private mgmtService: ManagementService,
    private newOrganizationService: NewOrganizationService,
  ) {
    effect(() => {
      const orgId = this.newOrganizationService.orgId();
      if (orgId) {
        this.orgId = orgId;
        this.loadProjectShortcuts();
      }
    });
  }

  public loadProjectShortcuts(): void {
    this.mgmtService.ownedProjects.pipe(takeUntil(this.destroy$)).subscribe((projects) => {
      if (projects) {
        const mapped: ShortcutItem[] = projects.map((p) => {
          const policy: ShortcutItem = {
            id: `project-${p.id}`,
            type: ShortcutType.PROJECT,
            title: p.name,
            i18nDesc: 'PROJECT.PAGES.TYPE.OWNED',
            routerLink: ['/projects', p.id],
            withRole: ['project.read', `project.read:${p.id}`],
            label: 'P',
            disabled: false,
            state: p.state,
          };
          return policy;
        });

        const routesShortcuts = [PROFILE_SHORTCUT, CREATE_ORG, CREATE_PROJECT, CREATE_USER];
        const settingsShortcuts = SETTINGLINKS.map((p) => {
          const policy: ShortcutItem = {
            id: p.i18nTitle,
            type: ShortcutType.POLICY,
            i18nTitle: p.i18nTitle,
            i18nDesc: p.i18nDesc,
            routerLink: p.orgRouterLink ?? p.iamRouterLink,
            queryParams: p.queryParams,
            withRole: p.orgWithRole ?? p.iamWithRole,
            icon: p.icon ?? '',
            svgIcon: p.svgIcon ?? '',
            color: p.color ?? '',
            disabled: false,
          };
          return policy;
        });

        this.ALL_SHORTCUTS = [...routesShortcuts, ...settingsShortcuts, ...mapped];
        this.loadShortcuts(this.orgId);
      }
    });
  }

  public loadShortcuts(orgId: string): void {
    ['main', 'secondary', 'third'].map((listName) => {
      const joinedShortcuts = this.storageService.getItem(`shortcuts:${listName}:${orgId}`, StorageLocation.local);
      if (joinedShortcuts) {
        const parsedIds: string[] = joinedShortcuts.split(',');
        if (parsedIds && parsedIds.length) {
          switch (listName) {
            case 'main':
              this.main = this.ALL_SHORTCUTS.filter((s) => parsedIds.includes(s.id));
              break;
            case 'secondary':
              this.secondary = this.ALL_SHORTCUTS.filter((s) => parsedIds.includes(s.id));
              break;
            case 'third':
              this.third = this.ALL_SHORTCUTS.filter((s) => parsedIds.includes(s.id));
              break;
          }
          this.organizeAllItems();
        }
      } else {
        switch (listName) {
          case 'main':
            this.main = [CREATE_ORG, CREATE_PROJECT, CREATE_USER];
            break;
          case 'secondary':
            this.secondary = this.ALL_SHORTCUTS.filter((item) => item.i18nTitle === 'SETTINGS.GROUPS.APPEARANCE');
            // [LOGIN_POLICY, PRIVATELABEL_POLICY].map((p) => {
            //   const policy: string = {
            //     i18nTitle: p.i18nTitle,
            //     i18nDesc: p.i18nDesc,
            //     routerLink: p.orgRouterLink,
            //     withRole: p.orgWithRole,
            //     icon: p.icon ?? '',
            //     color: p.color ?? '',
            //     disabled: false,
            //   };
            //   return policy;
            // });
            break;
          case 'third':
            this.third = [PROFILE_SHORTCUT];
            // [LOGIN_TEXTS_POLICY, MESSAGE_TEXTS_POLICY].map((p) => {
            //   const policy: ShortcutItem = {
            //     i18nTitle: p.i18nTitle,
            //     i18nDesc: p.i18nDesc,
            //     routerLink: p.orgRouterLink,
            //     withRole: p.orgWithRole,
            //     icon: p.icon ?? '',
            //     color: p.color ?? '',
            //     disabled: false,
            //   };
            //   return policy;
            // });
            break;
        }
        this.organizeAllItems();
      }
    });
  }

  private organizeAllItems(): void {
    const list = [this.main, this.secondary, this.third].flat();
    const filteredPolicies = this.ALL_SHORTCUTS.filter((p) => !list.find((l) => l.id === p.id));
    this.all = filteredPolicies;

    this.main === this.main.filter((s) => s.id && this.ALL_SHORTCUTS.map((p) => p.id).includes(s.id));
    this.secondary === this.secondary.filter((s) => s.id && this.ALL_SHORTCUTS.map((p) => p.id).includes(s.id));
    this.third === this.third.filter((s) => s.id && this.ALL_SHORTCUTS.map((p) => p.id).includes(s.id));
  }

  public ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public drop(event: CdkDragDrop<ShortcutItem[]>, listName: string) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
      this.saveStateToStorage();
      this.organizeAllItems();
    } else {
      transferArrayItem(event.previousContainer.data, event.container.data, event.previousIndex, event.currentIndex);
      this.saveStateToStorage();
      this.organizeAllItems();
    }
  }

  public saveStateToStorage(): void {
    const orgId = this.newOrganizationService.orgId();
    if (orgId) {
      this.storageService.setItem(`shortcuts:main:${orgId}`, this.main.map((p) => p.id).join(','), StorageLocation.local);
      this.storageService.setItem(
        `shortcuts:secondary:${orgId}`,
        this.secondary.map((p) => p.id).join(','),
        StorageLocation.local,
      );
      this.storageService.setItem(`shortcuts:third:${orgId}`, this.third.map((p) => p.id).join(','), StorageLocation.local);
    }
  }

  public reset(): void {
    const orgId = this.newOrganizationService.orgId();
    if (orgId) {
      ['main', 'secondary', 'third'].map((listName) => {
        this.storageService.removeItem(`shortcuts:${listName}:${orgId}`, StorageLocation.local);
      });

      this.loadShortcuts(orgId);
    }
  }

  public get allRoutes(): ShortcutItem[] {
    return this.all.filter((s) => s.type === ShortcutType.ROUTE);
  }

  public get allPolicies(): ShortcutItem[] {
    return this.all.filter((s) => s.type === ShortcutType.POLICY);
  }

  public get allProjects(): ShortcutItem[] {
    return this.all.filter((s) => s.type === ShortcutType.PROJECT);
  }

  public get allAvailableShortcuts(): ShortcutItem[] {
    return [...this.allRoutes, ...this.allPolicies, ...this.allProjects];
  }
}
