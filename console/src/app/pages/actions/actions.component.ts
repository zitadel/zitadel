import { Component } from '@angular/core';
import { enterAnimations } from 'src/app/animations';
import { SidenavSetting } from 'src/app/modules/sidenav/sidenav.component';
import { Breadcrumb, BreadcrumbService, BreadcrumbType } from 'src/app/services/breadcrumb.service';

const ACTIONS: SidenavSetting = { id: 'general', i18nKey: 'MENU.ACTIONS' };
const TARGETS: SidenavSetting = { id: 'roles', i18nKey: 'MENU.TARGETS' };

@Component({
  selector: 'cnsl-actions',
  templateUrl: './actions.component.html',
  styleUrls: ['./actions.component.scss'],
  animations: [enterAnimations],
})
export class ActionsComponent {
  public settingsList: SidenavSetting[] = [ACTIONS, TARGETS];
  public currentSetting = this.settingsList[0];

  constructor(breadcrumbService: BreadcrumbService) {
    const iamBread = new Breadcrumb({
      type: BreadcrumbType.INSTANCE,
      name: 'Instance',
      routerLink: ['/instance'],
    });

    breadcrumbService.setBreadcrumb([iamBread]);
  }
}
