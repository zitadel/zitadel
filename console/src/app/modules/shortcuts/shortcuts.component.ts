import { CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { Component, OnDestroy } from '@angular/core';
import { Subject, takeUntil } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
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
  label?: string;
  svgIcon?: string;
  avatarSrc?: string;
  color?: string;
  disabled?: boolean;
  state?: ProjectState;
}

const PROFILE_SHORTCUT: ShortcutItem = {
  routerLink: ['/users', 'me'],
  i18nTitle: 'USER.TITLE',
  icon: 'las la-cog',
  withRole: [''],
  disabled: false,
};

const CREATE_ORG: ShortcutItem = {
  i18nTitle: 'ORG.PAGES.CREATE',
  routerLink: ['/org', 'create'],
  withRole: ['org.create', 'iam.write'],
  icon: 'las la-plus',
  disabled: false,
};

const CREATE_PROJECT: ShortcutItem = {
  i18nTitle: 'PROJECT.PAGES.CREATE',
  routerLink: ['/projects', 'create'],
  withRole: ['project.create'],
  icon: 'las la-plus',
  disabled: false,
};

const CREATE_USER: ShortcutItem = {
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
})
export class ShortcutsComponent implements OnDestroy {
  public org!: Org.AsObject;
  public main: ShortcutItem[] = [];
  public secondary: ShortcutItem[] = [];
  public third: ShortcutItem[] = [];

  public ALL_ROUTES = [PROFILE_SHORTCUT, CREATE_ORG, CREATE_PROJECT, CREATE_USER];

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

  public ALL_PROJECTS: ShortcutItem[] = [];

  public allPolicies: ShortcutItem[] = this.ALL_POLICIES;
  public allProjects: ShortcutItem[] = this.ALL_PROJECTS;
  public allRoutes: ShortcutItem[] = this.ALL_ROUTES;

  private destroy$: Subject<void> = new Subject();
  public editState: boolean = false;
  public ProjectState: any = ProjectState;
  constructor(
    private storageService: StorageService,
    private auth: GrpcAuthService,
    private mgmtService: ManagementService,
  ) {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      this.org = org;
      this.loadShortcuts(org);
    }

    this.auth.activeOrgChanged.pipe(takeUntil(this.destroy$)).subscribe((org) => {
      this.loadShortcuts(org);
    });

    this.loadProjectShortcuts();
  }

  public loadProjectShortcuts(): void {
    this.mgmtService.ownedProjects.pipe(takeUntil(this.destroy$)).subscribe((projects) => {
      if (projects) {
        const mapped: ShortcutItem[] = projects.map((p) => {
          const policy: ShortcutItem = {
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

        this.ALL_PROJECTS = mapped;
        this.allProjects = this.ALL_PROJECTS;

        this.loadShortcuts(this.org);
      }
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
              this.main = parsed;
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
            this.main = [PROFILE_SHORTCUT, CREATE_ORG, CREATE_PROJECT, CREATE_USER];
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
    const filteredPolicies = this.ALL_POLICIES.filter((p) => !list.find((l) => l.i18nTitle === p.i18nTitle));
    this.allPolicies = filteredPolicies;

    const filteredProjects = this.ALL_PROJECTS.filter((p) => !list.find((l) => l.title === p.title));
    this.allProjects = filteredProjects;

    const filteredRoutes = this.ALL_ROUTES.filter((p) => !list.find((l) => l.i18nTitle === p.i18nTitle));
    this.allRoutes = filteredRoutes;

    this.main === this.main.filter((s) => s.title && this.ALL_PROJECTS.map((p) => p.title).includes(s.title));
    this.secondary === this.secondary.filter((s) => s.title && this.ALL_PROJECTS.map((p) => p.title).includes(s.title));
    this.third === this.third.filter((s) => s.title && this.ALL_PROJECTS.map((p) => p.title).includes(s.title));
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
