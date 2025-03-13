import { ChangeDetectionStrategy, Component, DestroyRef, signal } from '@angular/core';
import { WebKeysService } from 'src/app/services/webkeys.service';
import { NewFeatureService } from 'src/app/services/new-feature.service';
import { defer, EMPTY, firstValueFrom, Observable, ObservedValueOf, of, shareReplay, Subject, switchMap } from 'rxjs';
import { catchError, map, startWith, withLatestFrom } from 'rxjs/operators';
import { ToastService } from 'src/app/services/toast.service';
// todo: use v2beta
import { GetWebKey, WebKeyState } from '@zitadel/proto/zitadel/resources/webkey/v3alpha/key_pb';
import {
  WebKeyECDSAConfig_ECDSACurve,
  WebKeyRSAConfig_RSABits,
  WebKeyRSAConfig_RSAHasher,
} from '@zitadel/proto/zitadel/resources/webkey/v3alpha/config_pb';
import { MessageInitShape } from '@bufbuild/protobuf';
import { CreateWebKeyRequestSchema } from '@zitadel/proto/zitadel/resources/webkey/v3alpha/webkey_service_pb';
import { OidcWebKeysCreateComponent } from './oidc-webkeys-create/oidc-webkeys-create.component';
import { TimestampToDatePipe } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe';
import { MatDialog } from '@angular/material/dialog';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

const CACHE_WARNING_MS = 5 * 60 * 1000; // 5 minutes

@Component({
  selector: 'cnsl-oidc-webkeys',
  templateUrl: './oidc-webkeys.component.html',
  styleUrls: ['./oidc-webkeys.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysComponent {
  protected readonly refresh = new Subject<true>();
  protected readonly webKeysEnabled$: Observable<boolean>;
  protected readonly webKeys$: Observable<GetWebKey[]>;
  protected readonly inactiveWebKeys$: Observable<GetWebKey[]>;
  protected readonly nextWebKeyCandidate$: Observable<GetWebKey | undefined>;

  protected readonly activateLoading = signal(false);
  protected readonly createLoading = signal(false);

  constructor(
    private readonly webKeysService: WebKeysService,
    private readonly featureService: NewFeatureService,
    private readonly toast: ToastService,
    private readonly timestampToDatePipe: TimestampToDatePipe,
    private readonly dialog: MatDialog,
    private readonly destroyRef: DestroyRef,
  ) {
    this.webKeysEnabled$ = this.getWebKeysEnabled().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    const webKeys$ = this.getWebKeys(this.webKeysEnabled$).pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.webKeys$ = webKeys$.pipe(map((webKeys) => webKeys.filter((webKey) => webKey.state !== WebKeyState.STATE_INACTIVE)));
    this.inactiveWebKeys$ = webKeys$.pipe(
      map((webKeys) => webKeys.filter((webKey) => webKey.state === WebKeyState.STATE_INACTIVE)),
    );

    this.nextWebKeyCandidate$ = this.getNextWebKeyCandidate(this.webKeys$);
  }

  private getWebKeysEnabled() {
    return defer(() => this.featureService.getInstanceFeatures()).pipe(
      map((features) => features.webKey?.enabled ?? false),
      catchError((err) => {
        this.toast.showError(err);
        return of(false);
      }),
    );
  }

  private getWebKeys(webKeysEnabled$: Observable<boolean>) {
    return this.refresh.pipe(
      startWith(true),
      withLatestFrom(webKeysEnabled$),
      switchMap(([_, enabled]) => {
        if (!enabled) {
          return EMPTY;
        }
        return this.webKeysService.ListWebKeys();
      }),
      map(({ webKeys }) => webKeys),
      catchError((err) => {
        this.toast.showError(err);
        return EMPTY;
      }),
    );
  }

  private getNextWebKeyCandidate(webKeys$: Observable<GetWebKey[]>) {
    return webKeys$.pipe(
      map((webKeys) => {
        if (webKeys.length < 2) {
          return undefined;
        }
        const [webKey, nextWebKey] = webKeys;
        if (webKey.state !== WebKeyState.STATE_ACTIVE) {
          return undefined;
        }
        if (nextWebKey.state !== WebKeyState.STATE_INITIAL) {
          return undefined;
        }
        return nextWebKey;
      }),
    );
  }

  protected async createWebKey(event: ObservedValueOf<OidcWebKeysCreateComponent['ngSubmit']>) {
    try {
      this.createLoading.set(true);

      const req = !event
        ? this.createEd25519()
        : 'curve' in event
          ? this.createEcdsa(event.curve)
          : this.createRsa(event.bits, event.hasher);
      await this.webKeysService.CreateWebKey(req);

      this.refresh.next(true);
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.createLoading.set(false);
    }
  }

  private createEd25519(): MessageInitShape<typeof CreateWebKeyRequestSchema> {
    return {
      key: {
        case: 'ed25519',
        value: {},
      },
    } as any;
  }

  private createEcdsa(curve: WebKeyECDSAConfig_ECDSACurve): MessageInitShape<typeof CreateWebKeyRequestSchema> {
    // todo: use correct typing
    return {
      key: {
        case: 'ecdsa',
        value: {
          curve,
        },
      },
    } as any;
  }

  private createRsa(
    bits: WebKeyRSAConfig_RSABits,
    hasher: WebKeyRSAConfig_RSAHasher,
  ): MessageInitShape<typeof CreateWebKeyRequestSchema> {
    // todo: use correct typing
    return {
      key: {
        case: 'rsa',
        value: {
          bits,
          hasher,
        },
      },
    } as any;
  }

  protected async deleteWebKey(row: GetWebKey) {
    // todo: fix this when typings are correct
    try {
      await this.webKeysService.DeleteWebKey((row as any).id);
      this.refresh.next(true);
    } catch (err) {
      this.toast.showError(err);
    }
  }

  protected async activateWebKey(nextWebKey: GetWebKey) {
    try {
      this.activateLoading.set(true);
      // todo: fix this when typing are correct
      const creationDate = this.timestampToDatePipe.transform((nextWebKey as any).creationDate);
      if (!creationDate) {
        // noinspection ExceptionCaughtLocallyJS
        throw new Error('Invalid creation date');
      }

      const diffToCurrentTime = Date.now() - creationDate.getTime();
      if (diffToCurrentTime < CACHE_WARNING_MS && !(await this.openCacheWarnDialog())) {
        return;
      }

      // todo: remove this once typing is fixed
      await this.webKeysService.ActivateWebKey((nextWebKey as any).id);
      this.refresh.next(true);
    } catch (error) {
      this.toast.showError(error);
    } finally {
      this.activateLoading.set(false);
    }
  }

  private openCacheWarnDialog() {
    const dialogRef = this.dialog.open(WarnDialogComponent, {
      data: {
        confirmKey: 'ACTIONS.DELETE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'IDP.DELETE_TITLE',
        descriptionKey: 'IDP.DELETE_DESCRIPTION',
      },
      width: '400px',
    });

    const obs = dialogRef.afterClosed().pipe(map(Boolean), takeUntilDestroyed(this.destroyRef));
    return firstValueFrom(obs);
  }
}
