import { ChangeDetectionStrategy, Component, DestroyRef, OnInit, signal } from '@angular/core';
import { WebKeysService } from 'src/app/services/webkeys.service';
import { defer, EMPTY, firstValueFrom, Observable, ObservedValueOf, of, shareReplay, Subject, switchMap } from 'rxjs';
import { catchError, map, startWith, withLatestFrom } from 'rxjs/operators';
import { ToastService } from 'src/app/services/toast.service';
import { MessageInitShape } from '@bufbuild/protobuf';
import { OidcWebKeysCreateComponent } from './oidc-webkeys-create/oidc-webkeys-create.component';
import { TimestampToDatePipe } from 'src/app/pipes/timestamp-to-date-pipe/timestamp-to-date.pipe';
import { MatDialog } from '@angular/material/dialog';
import { WarnDialogComponent } from 'src/app/modules/warn-dialog/warn-dialog.component';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { State, WebKey } from '@zitadel/proto/zitadel/webkey/v2beta/key_pb';
import { CreateWebKeyRequestSchema } from '@zitadel/proto/zitadel/webkey/v2beta/webkey_service_pb';
import { RSAHasher, RSABits, ECDSACurve } from '@zitadel/proto/zitadel/webkey/v2beta/key_pb';

const CACHE_WARNING_MS = 5 * 60 * 1000; // 5 minutes

@Component({
  selector: 'cnsl-oidc-webkeys',
  templateUrl: './oidc-webkeys.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OidcWebKeysComponent {
  protected readonly refresh = new Subject<true>();
  protected readonly webKeys$: Observable<WebKey[]>;
  protected readonly inactiveWebKeys$: Observable<WebKey[]>;
  protected readonly nextWebKeyCandidate$: Observable<WebKey | undefined>;

  protected readonly activateLoading = signal(false);
  protected readonly createLoading = signal(false);

  constructor(
    private readonly webKeysService: WebKeysService,
    private readonly toast: ToastService,
    private readonly timestampToDatePipe: TimestampToDatePipe,
    private readonly dialog: MatDialog,
    private readonly destroyRef: DestroyRef,
  ) {
    const webKeys$ = this.getWebKeys().pipe(shareReplay({ refCount: true, bufferSize: 1 }));

    this.webKeys$ = webKeys$.pipe(map((webKeys) => webKeys.filter((webKey) => webKey.state !== State.INACTIVE)));
    this.inactiveWebKeys$ = webKeys$.pipe(map((webKeys) => webKeys.filter((webKey) => webKey.state === State.INACTIVE)));

    this.nextWebKeyCandidate$ = this.getNextWebKeyCandidate(this.webKeys$);
  }

  private getWebKeys() {
    return this.refresh.pipe(
      startWith(true),
      switchMap(() => {
        return this.webKeysService.ListWebKeys();
      }),
      map(({ webKeys }) => webKeys),
      catchError(async (err) => {
        this.toast.showError(err);
        return [];
      }),
    );
  }

  private getNextWebKeyCandidate(webKeys$: Observable<WebKey[]>) {
    return webKeys$.pipe(
      map((webKeys) => {
        if (webKeys.length < 2) {
          return undefined;
        }
        const [webKey, nextWebKey] = webKeys;
        if (webKey.state !== State.ACTIVE) {
          return undefined;
        }
        if (nextWebKey.state !== State.INITIAL) {
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
    };
  }

  private createEcdsa(curve: ECDSACurve): MessageInitShape<typeof CreateWebKeyRequestSchema> {
    return {
      key: {
        case: 'ecdsa',
        value: {
          curve,
        },
      },
    };
  }

  private createRsa(bits: RSABits, hasher: RSAHasher): MessageInitShape<typeof CreateWebKeyRequestSchema> {
    return {
      key: {
        case: 'rsa',
        value: {
          bits,
          hasher,
        },
      },
    };
  }

  protected async deleteWebKey(row: WebKey) {
    try {
      await this.webKeysService.DeleteWebKey(row.id);
      this.refresh.next(true);
    } catch (err) {
      this.toast.showError(err);
    }
  }

  protected async activateWebKey(nextWebKey: WebKey) {
    try {
      this.activateLoading.set(true);
      const creationDate = this.timestampToDatePipe.transform(nextWebKey.creationDate);
      if (!creationDate) {
        // noinspection ExceptionCaughtLocallyJS
        throw new Error('Invalid creation date');
      }

      const diffToCurrentTime = Date.now() - creationDate.getTime();
      if (diffToCurrentTime < CACHE_WARNING_MS && !(await this.openCacheWarnDialog())) {
        return;
      }

      await this.webKeysService.ActivateWebKey(nextWebKey.id);
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
        confirmKey: 'DESCRIPTIONS.SETTINGS.WEB_KEYS.TABLE.ACTIVATE',
        cancelKey: 'ACTIONS.CANCEL',
        titleKey: 'Web Key is less then 5 min old',
        descriptionKey: 'DESCRIPTIONS.SETTINGS.WEB_KEYS.TABLE.NOTE',
      },
      width: '400px',
    });

    const obs = dialogRef.afterClosed().pipe(map(Boolean), takeUntilDestroyed(this.destroyRef));
    return firstValueFrom(obs);
  }
}
