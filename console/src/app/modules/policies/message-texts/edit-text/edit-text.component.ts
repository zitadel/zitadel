import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'cnsl-edit-text',
  templateUrl: './edit-text.component.html',
  styleUrls: ['./edit-text.component.scss']
})
export class EditTextComponent implements OnInit {
  @Input() label: string = 'hello';
  constructor() { }

  ngOnInit(): void {
  }

}
