import { ChangeDetectorRef, Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';

import { Theme, View } from '../private-labeling-policy.component';

declare const tinycolor: any;

@Component({
  selector: 'cnsl-preview',
  templateUrl: './preview.component.html',
  styleUrls: ['./preview.component.scss'],
})
export class PreviewComponent implements OnInit, OnDestroy {
  @Input() policy!: LabelPolicy.AsObject;
  @Input() theme: Theme = Theme.DARK;
  @Input() refresh: Observable<void> = of();
  private destroyed$: Subject<void> = new Subject();
  public Theme: any = Theme;
  public View: any = View;
  constructor(private chd: ChangeDetectorRef) {}

  public ngOnInit(): void {
    this.refresh.pipe(takeUntil(this.destroyed$)).subscribe(() => {
      this.chd.detectChanges();
    });
  }

  public ngOnDestroy(): void {
    this.destroyed$.next();
    this.destroyed$.complete();
  }

  public get contrastTextColor(): string {
    const c = tinycolor(this.theme === Theme.DARK ? this.policy.primaryColorDark : this.policy.primaryColor);
    return this.getContrast(c.toHexString());
  }

  public getContrast(color: string): string {
    const onBlack = tinycolor.readability('#000', color);
    const onWhite = tinycolor.readability('#fff', color);
    if (onBlack > onWhite) {
      return '#000000';
    } else {
      return '#ffffff';
    }
  }
}
