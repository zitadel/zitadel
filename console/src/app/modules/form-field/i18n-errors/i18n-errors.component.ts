import { KeyValue, KeyValuePipe } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { distinctUntilChanged, filter, map, mergeMap, Observable } from 'rxjs';

@Component({
  selector: 'cnsl-i18n-errors',
  templateUrl: './i18n-errors.component.html',
})
export class I18nErrorsComponent implements OnInit {
  @Input() ctrl!: AbstractControl;

  public errori18nKeys$!: Observable<Array<string>>;

  constructor(    private kvPipe: KeyValuePipe    ) {}

    ngOnInit(): void {
      this.errori18nKeys$ = this.ctrl.valueChanges.pipe(
        mergeMap(() => this.ctrl.statusChanges),
        map(() => this.currentErrors()),
      )
    }

  private currentErrors(): Array<string> {
    return this.kvPipe.transform(this.ctrl.errors)?.
    filter(this.filterErrorsProperties).
    map(this.mapErrorToI18nKey).
    filter(this.distinctFilter) || []
  }

  private filterErrorsProperties(kv: KeyValue<unknown,unknown>): boolean {
    console.log("filterErrorsProperties", kv)
    return (kv.key as string) != "invalid" &&
      kv.key as string != "required" &&
      !(kv.value as any).invalid
  }

  private mapErrorToI18nKey(kv: KeyValue<unknown,unknown>): string {
    console.log("mapErrorToI18nKey", kv)
    return (kv.value as {i18nKey: string}).i18nKey || "ERRORS.INVALID_FORMAT"
  }

  private distinctFilter(item: string, index: number, arr: Array<string>): boolean{
    return arr.indexOf(item) === index
  }
}
