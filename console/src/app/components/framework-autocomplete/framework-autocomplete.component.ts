import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output, signal } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import frameworkDefinition from '../../../../../docs/frameworks.json';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { MatAutocompleteModule, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { InputModule } from 'src/app/modules/input/input.module';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { Observable, map, startWith, switchMap, tap } from 'rxjs';

type FrameworkDefinition = {
  title: string;
  imgSrcDark: string;
  imgSrcLight?: string;
  docsLink: string;
  external?: boolean;
};

type Framework = FrameworkDefinition & {
  fragment: string;
};

@Component({
  standalone: true,
  selector: 'cnsl-framework-autocomplete',
  templateUrl: './framework-autocomplete.component.html',
  styleUrls: ['./framework-autocomplete.component.scss'],
  imports: [
    TranslateModule,
    RouterModule,
    MatSelectModule,
    MatAutocompleteModule,
    ReactiveFormsModule,
    MatProgressSpinnerModule,
    FormsModule,
    CommonModule,
    MatButtonModule,
    InputModule,
  ],
})
export class FrameworkAutocompleteComponent {
  public isLoading = signal(false);
  public frameworks: Framework[] = frameworkDefinition.map((f) => {
    return {
      ...f,
      fragment: '',
      imgSrcDark: `assets${f.imgSrcDark}`,
      imgSrcLight: `assets${f.imgSrcLight ? f.imgSrcLight : f.imgSrcDark}`,
    };
  });
  public myControl: FormControl = new FormControl();
  @Output() public selectionChanged: EventEmitter<Framework> = new EventEmitter();
  filteredOptions: Observable<Framework[]>;

  constructor() {
    this.filteredOptions = this.myControl.valueChanges.pipe(
      startWith(''),
      map((value) => this._filter(value || '')),
    );
  }

  private _filter(value: string): Framework[] {
    const filterValue = value.toLowerCase();

    return this.frameworks.filter((option) => option.title.toLowerCase().includes(filterValue));
  }

  public displayFn(project?: any): string {
    return project && project.projectName ? `${project.projectName}` : project && project.name ? `${project.name}` : '';
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    const f: Framework = event.option.value;
    this.selectionChanged.emit(f);
  }
}
