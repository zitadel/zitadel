import { Directive, Input } from '@angular/core';
import { DataSource } from '@angular/cdk/collections';
import { MatCellDef } from '@angular/material/table';
import { CdkCellDef } from '@angular/cdk/table';

@Directive({
  selector: '[cnslCellDef]',
  providers: [{ provide: CdkCellDef, useExisting: TypeSafeCellDefDirective }],
})
export class TypeSafeCellDefDirective<T> extends MatCellDef {
  @Input({ required: true }) cnslCellDefDataSource!: DataSource<T>;

  static ngTemplateContextGuard<T>(_dir: TypeSafeCellDefDirective<T>, _ctx: any): _ctx is { $implicit: T; index: number } {
    return true;
  }
}
