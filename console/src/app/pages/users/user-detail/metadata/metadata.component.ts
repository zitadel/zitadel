import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';

import { MetadataDialogComponent } from '../metadata-dialog/metadata-dialog.component';

@Component({
  selector: 'cnsl-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.scss'],
})
export class MetadataComponent implements OnInit {
  @Input() userId: string = '';
  @Input() serviceType: string = '';
  private service!: GrpcAuthService | ManagementService;
  public metadata: Metadata.AsObject[] = [];


  constructor(private dialog: MatDialog, private injector: Injector,
  ) { }

  ngOnInit(): void {
    if (this.userId) {
      this.service = this.injector.get(ManagementService as Type<ManagementService>);
    } else {
      this.service = this.injector.get(GrpcAuthService as Type<GrpcAuthService>);
    }

    this.loadMetadata();
  }

  public editMetadata(): void {
    const dialogRef = this.dialog.open(MetadataDialogComponent, {
      data: {
        serviceType: this.userId ? 'MGMT' : 'AUTH',
        userId: this.userId,
      },
    });

    dialogRef.afterClosed().subscribe(resp => {
      if (resp) {

      }
    });
  }

  public loadMetadata(userId?: string): void {
    if (userId && this.serviceType === 'MGMT') {
      (this.service as ManagementService).listUserMetadata(userId).then(resp => {
        this.metadata = resp.resultList;
      });
    } else if (this.serviceType === 'AUTH') {
      (this.service as GrpcAuthService).listMyMetadata().then(resp => {
        this.metadata = resp.resultList;
        console.log(this.metadata);
      });
    }
  }
}
