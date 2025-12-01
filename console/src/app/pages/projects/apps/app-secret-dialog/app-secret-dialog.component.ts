import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import JSZip from 'jszip';
import { saveAs } from 'file-saver';
import { EnvironmentService } from 'src/app/services/environment.service';
import { take } from 'rxjs/operators';

@Component({
  selector: 'cnsl-app-secret-dialog',
  templateUrl: './app-secret-dialog.component.html',
  styleUrls: ['./app-secret-dialog.component.scss'],
  standalone: false,
})
export class AppSecretDialogComponent {
  public copied: string = '';
  constructor(
    public dialogRef: MatDialogRef<AppSecretDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
    private envService: EnvironmentService,
  ) { }

  public closeDialog(): void {
    this.dialogRef.close(false);
  }

  public async downloadTemplate() {
    const zip = new JSZip();
    const response = await fetch(`assets/templates/${this.data.frameworkId}.zip`);
    const blob = await response.blob();
    const templateZip = await zip.loadAsync(blob);

    const env = await this.envService.env.pipe(take(1)).toPromise();
    if (!env) {
      return;
    }

    const envFile = `VITE_ZITADEL_ISSUER=${env.issuer}
VITE_ZITADEL_CLIENT_ID=${this.data.clientId}`;

    templateZip.file('.env', envFile);

    const content = await templateZip.generateAsync({ type: 'blob' });
    saveAs(content, 'zitadel-react-app.zip');
  }
}
