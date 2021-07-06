import { Component, Input, OnInit } from '@angular/core';
import { Observable } from 'rxjs';

@Component({
  selector: 'cnsl-edit-text',
  templateUrl: './edit-text.component.html',
  styleUrls: ['./edit-text.component.scss']
})
export class EditTextComponent implements OnInit {
  @Input() label: string = 'hello';
  @Input() current$: Observable<string>;

  private value: string = '';
  constructor() { }

  ngOnInit(): void {
  }

}
