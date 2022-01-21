import { Component, EventEmitter, HostListener, Input, Output } from '@angular/core';

export enum ActionKeysType {
  ADD,
  DELETE,
  DEACTIVATE,
  REACTIVATE,
}

@Component({
  selector: 'cnsl-action-keys',
  templateUrl: './action-keys.component.html',
  styleUrls: ['./action-keys.component.scss'],
})
export class ActionKeysComponent {
  @Input() type: ActionKeysType = ActionKeysType.ADD;
  @Input() withoutMargin: boolean = false;
  @Input() doNotUseContrast: boolean = false;
  @Output() actionTriggered: EventEmitter<void> = new EventEmitter();
  @HostListener('document:keydown', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent) {
    switch (this.type) {
      case ActionKeysType.ADD:
        if (event.code === 'KeyN') {
          this.actionTriggered.emit();
        }
        break;

      case ActionKeysType.DELETE:
        if ((event.ctrlKey || event.metaKey) && event.code === 'Enter') {
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
    }
  }
  public ActionKeysType: any = ActionKeysType;
  constructor() {}

  public get isMacLike(): boolean {
    return /(Mac|iPhone|iPod|iPad)/i.test(navigator.userAgent);
  }

  public get isIOS(): boolean {
    return /(iPhone|iPod|iPad)/i.test(navigator.userAgent);
  }
}
