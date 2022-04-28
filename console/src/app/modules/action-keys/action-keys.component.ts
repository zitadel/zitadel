import { BreakpointObserver, Breakpoints } from '@angular/cdk/layout';
import { AfterViewInit, Component, EventEmitter, HostListener, Input, Output } from '@angular/core';
import { map, Observable } from 'rxjs';

export enum ActionKeysType {
  ADD,
  DELETE,
  DEACTIVATE,
  REACTIVATE,
  FILTER,
  ORGSWITCHER,
  CLEAR,
}

@Component({
  selector: 'cnsl-action-keys',
  templateUrl: './action-keys.component.html',
  styleUrls: ['./action-keys.component.scss'],
})
export class ActionKeysComponent implements AfterViewInit {
  @Input() type: ActionKeysType = ActionKeysType.ADD;
  @Input() withoutMargin: boolean = false;
  @Input() doNotUseContrast: boolean = false;
  @Output() actionTriggered: EventEmitter<void> = new EventEmitter();
  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent) {
    const tagname = (event.target as any)?.tagName;
    const exclude = ['input', 'textarea'];

    if (exclude.indexOf(tagname.toLowerCase()) === -1) {
      switch (this.type) {
        case ActionKeysType.CLEAR:
          if (event.code === 'Escape') {
            event.preventDefault();
            this.actionTriggered.emit();
          }
          break;
        case ActionKeysType.ORGSWITCHER:
          if (event.key === '/') {
            this.actionTriggered.emit();
          }
          break;
        case ActionKeysType.ADD:
          if (event.code === 'KeyN') {
            this.actionTriggered.emit();
          }
          break;

        case ActionKeysType.DELETE:
          if ((event.ctrlKey || event.metaKey) && event.code === 'Backspace') {
            this.actionTriggered.emit();
          }
          break;

        case ActionKeysType.DEACTIVATE:
          if ((event.ctrlKey || event.metaKey) && event.code === 'ArrowDown') {
            event.preventDefault();
            this.actionTriggered.emit();
          }
          break;

        case ActionKeysType.REACTIVATE:
          if ((event.ctrlKey || event.metaKey) && event.code === 'ArrowUp') {
            event.preventDefault();
            this.actionTriggered.emit();
          }
          break;

        case ActionKeysType.FILTER:
          if (event.ctrlKey === false && event.code === 'KeyF') {
            this.actionTriggered.emit();
          }
          break;
      }
    }
  }
  public isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset).pipe(
    map((result) => {
      return result.matches;
    }),
  );

  public ActionKeysType: any = ActionKeysType;

  constructor(public breakpointObserver: BreakpointObserver) {}

  ngAfterViewInit(): void {
    window.focus();
    if (document.activeElement) {
      (document.activeElement as any).blur();
    }
  }

  public get isMacLike(): boolean {
    return /(Mac|iPhone|iPod|iPad)/i.test(navigator.userAgent);
  }

  public get isIOS(): boolean {
    return /(iPhone|iPod|iPad)/i.test(navigator.userAgent);
  }
}
