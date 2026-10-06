import { CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { Component, OnDestroy } from '@angular/core';
import { combineLatest, EMPTY, from, map, startWith, Subject, switchMap, takeUntil } from 'rxjs';
import { Org } from 'src/app/proto/generated/zitadel/org_pb';
import { Project, ProjectState } from 'src/app/proto/generated/zitadel/project_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { StorageLocation, StorageService } from 'src/app/services/storage.service';

import { SETTINGLINKS } from '../settings-grid/settinglinks';

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

const PROJECT_SHORTCUT_PREFIX = 'project-';
const SHORTCUT_LISTS = ['main', 'secondary', 'third'];

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
  public org!: Org.AsObject;

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
    private auth: GrpcAuthService,
    private mgmtService: ManagementService,
  ) {
    combineLatest([this.auth.activeOrgChanged.pipe(startWith(undefined)), this.mgmtService.ownedProjects])
      .pipe(
        switchMap(([, projects]) => {
          const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
          if (!org?.id) {
            return EMPTY;
          }
          return from(this.withPinnedProjects(org, projects)).pipe(map((allProjects) => ({ org, allProjects })));
        }),
        takeUntil(this.destroy$),
      )
      .subscribe(({ org, allProjects }) => {
        this.org = org;
        this.loadProjectShortcuts(allProjects);
      });
  }

  public loadProjectShortcuts(projects: Project.AsObject[]): void {
    const mapped: ShortcutItem[] = projects.map((p) => {
      const policy: ShortcutItem = {
        id: `${PROJECT_SHORTCUT_PREFIX}${p.id}`,
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
    this.loadShortcuts(this.org);
  }

  // only a part of the projects is preloaded, so pinned projects might need to be loaded
  private async withPinnedProjects(org: Org.AsObject, projects: Project.AsObject[]): Promise<Project.AsObject[]> {
    const loadedIds = new Set(projects.map((p) => p.id));
    const missingIds = SHORTCUT_LISTS.flatMap((listName) =>
      (this.storageService.getItem(`shortcuts:${listName}:${org.id}`, StorageLocation.local) ?? '').split(','),
    )
      .filter((id: string) => id.startsWith(PROJECT_SHORTCUT_PREFIX))
      .map((id: string) => id.slice(PROJECT_SHORTCUT_PREFIX.length))
      .filter((id: string) => !loadedIds.has(id));
    if (!missingIds.length) {
      return projects;
    }

    const pinned = await Promise.all(
      [...new Set(missingIds)].map((id) =>
        this.mgmtService
          .getProjectByID(id)
          .then((resp) => resp.project)
          // e.g. the project was removed
          .catch(() => undefined),
      ),
    );
    return [...projects, ...pinned.filter((p): p is Project.AsObject => !!p)];
  }

  public loadShortcuts(org: Org.AsObject): void {
    SHORTCUT_LISTS.map((listName) => {
      const joinedShortcuts = this.storageService.getItem(`shortcuts:${listName}:${org.id}`, StorageLocation.local);
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

  public drop(event: CdkDragDrop<ShortcutItem[]>) {
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
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      this.storageService.setItem(`shortcuts:main:${org.id}`, this.main.map((p) => p.id).join(','), StorageLocation.local);
      this.storageService.setItem(
        `shortcuts:secondary:${org.id}`,
        this.secondary.map((p) => p.id).join(','),
        StorageLocation.local,
      );
      this.storageService.setItem(`shortcuts:third:${org.id}`, this.third.map((p) => p.id).join(','), StorageLocation.local);
    }
  }

  public reset(): void {
    const org: Org.AsObject | null = this.storageService.getItem('organization', StorageLocation.session);
    if (org && org.id) {
      SHORTCUT_LISTS.map((listName) => {
        this.storageService.removeItem(`shortcuts:${listName}:${org.id}`, StorageLocation.local);
      });

      this.loadShortcuts(org);
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
