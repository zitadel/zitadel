import { Component, signal } from '@angular/core';
import { enterAnimations } from 'src/app/animations';
import { InfoSectionType } from 'src/app/modules/info-section/info-section.component';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

const ACTIONS: SidenavSetting = { id: 'actions', i18nKey: 'MENU.ACTIONS' };
const TARGETS: SidenavSetting = { id: 'targets', i18nKey: 'MENU.TARGETS' };

@Component({
  selector: 'cnsl-actions',
  templateUrl: './actions.component.html',
  styleUrls: ['./actions.component.scss'],
  animations: [enterAnimations],
  standalone: false,
})
export class ActionsComponent {
  public settingsList: SidenavSetting[] = [ACTIONS, TARGETS];
  protected readonly currentSetting$ = signal<SidenavSetting>(this.settingsList[0]);
  protected readonly InfoSectionType = InfoSectionType;

  constructor(breadcrumbService: BreadcrumbService) {
    const iamBread = new Breadcrumb({
      type: BreadcrumbType.INSTANCE,
      name: 'Instance',
      routerLink: ['/instance'],
    });

    breadcrumbService.setBreadcrumb([iamBread]);
  }
}
