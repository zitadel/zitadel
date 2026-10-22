import { AfterViewChecked, Directive, DoCheck, ElementRef, HostListener } from '@angular/core';
import { MatTooltip } from '@angular/material/tooltip';
import { TranslateService } from '@ngx-translate/core';
import { CardComponent } from 'src/app/modules/card/card.component';

/** Adds semantics to native tooltip-only controls without replacing their templates/actions. */
// eslint-disable-next-line @angular-eslint/directive-selector -- Match unchanged native templates only within the Account compilation scope.
@Directive({ selector: 'button[matTooltip], a[matTooltip]', standalone: false })
export class AccountNativeActions implements DoCheck {
  constructor(
    private element: ElementRef<HTMLElement>,
    private tooltip: MatTooltip,
  ) {}
  ngDoCheck(): void {
    const el = this.element.nativeElement;
    if (this.tooltip.message && el.hasAttribute('mat-icon-button')) el.setAttribute('aria-label', this.tooltip.message);
    if (el.tagName === 'A' && !el.hasAttribute('href')) {
      el.setAttribute('role', 'button');
      el.setAttribute('tabindex', '0');
    }
  }
  @HostListener('keydown', ['$event'])
  activate(event: Event): void {
    const el = this.element.nativeElement;
    if (
      event instanceof KeyboardEvent &&
      el.tagName === 'A' &&
      !el.hasAttribute('href') &&
      ['Enter', ' '].includes(event.key)
    ) {
      event.preventDefault();
      el.click();
    }
  }
}

/** The shared card lives in its own NgModule; adapt its collapse button at the Account boundary. */
// eslint-disable-next-line @angular-eslint/directive-selector -- Adapter must attach to the unchanged shared component.
@Directive({ selector: 'cnsl-card', standalone: false })
export class AccountNativeCard implements AfterViewChecked {
  constructor(
    private element: ElementRef<HTMLElement>,
    private card: CardComponent,
    private translate: TranslateService,
  ) {}
  ngAfterViewChecked(): void {
    const button = this.element.nativeElement.querySelector<HTMLButtonElement>('.header .row > button.button');
    if (!button) return;
    button.setAttribute('aria-label', this.translate.instant(this.card.expanded ? 'ACCOUNT.COLLAPSE' : 'ACCOUNT.EXPAND'));
    button.setAttribute('aria-expanded', String(this.card.expanded));
  }
}

/** Associate shared form labels with Material comboboxes (the native label-for alone does not name them). */
// eslint-disable-next-line @angular-eslint/directive-selector -- Adapter must attach to the unchanged shared form.
@Directive({ selector: 'cnsl-detail-form, cnsl-paginator, cnsl-form-field', standalone: false })
export class AccountNativeFields implements AfterViewChecked {
  constructor(
    private element: ElementRef<HTMLElement>,
    private translate: TranslateService,
  ) {}
  ngAfterViewChecked(): void {
    const root = this.element.nativeElement;
    if (root.matches('cnsl-detail-form')) {
      const edit = root.querySelector<HTMLButtonElement>('.usernamediv button.edit');
      if (edit && !edit.querySelector('.la-pen')) {
        const label = edit.textContent?.trim() || this.translate.instant('USER.PROFILE.CHANGEUSERNAME');
        edit.setAttribute('aria-label', label);
        edit.title = label;
        const icon = document.createElement('i');
        icon.className = 'las la-pen';
        icon.setAttribute('aria-hidden', 'true');
        edit.replaceChildren(icon);
      }
    }
    const fields = root.matches('cnsl-form-field') ? [root] : root.querySelectorAll('cnsl-form-field');
    for (const field of fields) {
      const label = field.querySelector('label');
      const select = field.querySelector('mat-select');
      if (!select?.id) continue;
      if (!label) {
        if (select.classList.contains('paginator-select'))
          select.setAttribute('aria-label', this.translate.instant('ACCOUNT.PAGE_SIZE'));
        continue;
      }
      label.id ||= `account-label-${select.id}`;
      select.setAttribute('aria-labelledby', label.id);
    }
  }
}
