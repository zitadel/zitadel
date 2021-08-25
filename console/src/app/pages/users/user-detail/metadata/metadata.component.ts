import { Component, Injector, Input, OnInit, Type } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Timestamp } from 'google-protobuf/google/protobuf/timestamp_pb';
import { Metadata } from 'src/app/proto/generated/zitadel/metadata_pb';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

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
  public ts!: Timestamp.AsObject | undefined;
  public loading: boolean = false;

  constructor(private dialog: MatDialog, private injector: Injector, private toast: ToastService,
  ) { }

  ngOnInit(): void {
    console.log(this.userId);
    if (this.serviceType === 'MGMT') {
      this.service = this.injector.get(ManagementService as Type<ManagementService>);
    } else if (this.serviceType === 'AUTH') {
      this.service = this.injector.get(GrpcAuthService as Type<GrpcAuthService>);
    }

    this.loadMetadata();
  }

  public editMetadata(): void {
    const dialogRef = this.dialog.open(MetadataDialogComponent, {
      data: {
        serviceType: this.serviceType,
        userId: this.userId,
      },
    });

    dialogRef.afterClosed().subscribe(() => {
      this.loadMetadata();
    });
  }

  public loadMetadata(): Promise<any> {
    this.loading = true;
    if (this.userId && this.serviceType === 'MGMT') {
      return (this.service as ManagementService).listUserMetadata(this.userId).then(resp => {
        this.loading = false;
        this.metadata = resp.resultList.map(md => {
          return {
            key: md.key,
            value: atob(md.value as string),
          };
        });
        this.ts = resp.details?.viewTimestamp;
      }).catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
    } else {
      return (this.service as GrpcAuthService).listMyMetadata().then(resp => {
        this.loading = false;
        this.metadata = resp.resultList.map(md => {
          return {
            key: md.key,
            value: atob(md.value as string),
          };
        });
        this.ts = resp.details?.viewTimestamp;
      }).catch((error) => {
        this.loading = false;
        this.toast.showError(error);
      });
    }
  }
}
