import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { AssetEndpoint, AssetService } from 'src/app/services/asset.service';
import { EnvironmentService } from 'src/app/services/environment.service';
import { StorageService } from 'src/app/services/storage.service';

// Below both the product's 500 kB cap and the native 512 KiB limit. Re-encoding strips photo metadata.
export const AVATAR_MAX_BYTES = 480000;
export async function compressAvatar(file: File): Promise<File> {
  if (!file.type.startsWith('image/') || file.size > 30 * 1024 * 1024) throw new Error('avatar');
  const bitmap = await createImageBitmap(file, { imageOrientation: 'from-image' });
  try {
    const ratio = Math.min(1, 1024 / Math.max(bitmap.width, bitmap.height));
    let width = Math.max(1, Math.round(bitmap.width * ratio));
    let height = Math.max(1, Math.round(bitmap.height * ratio));
    const canvas = document.createElement('canvas');
    for (let attempt = 0; attempt < 6; attempt++) {
      canvas.width = width;
      canvas.height = height;
      const context = canvas.getContext('2d');
      if (!context) throw new Error('canvas');
      context.fillStyle = '#fff';
      context.fillRect(0, 0, width, height);
      context.drawImage(bitmap, 0, 0, width, height);
      for (const quality of [0.9, 0.8, 0.65, 0.5]) {
        const blob = await new Promise<Blob | null>((resolve) => canvas.toBlob(resolve, 'image/jpeg', quality));
        if (blob && blob.size <= AVATAR_MAX_BYTES) return new File([blob], 'avatar.jpg', { type: 'image/jpeg' });
      }
      width = Math.max(1, Math.round(width * 0.75));
      height = Math.max(1, Math.round(height * 0.75));
    }
    throw new Error('avatar');
  } finally {
    bitmap.close();
  }
}

/** Preserve the native upload/refresh/dialog flow; compress only the self-avatar payload. */
@Injectable()
export class AccountAssetService extends AssetService {
  constructor(
    env: EnvironmentService,
    http: HttpClient,
    storage: StorageService,
    private translate: TranslateService,
  ) {
    super(env, http, storage);
  }
  override async upload(endpoint: AssetEndpoint | string, body: FormData, orgId?: string): Promise<unknown> {
    if (endpoint !== 'users/me/avatar') return super.upload(endpoint, body, orgId);
    const file = body.get('file');
    try {
      if (!(file instanceof File)) throw new Error('avatar');
      const compressed = await compressAvatar(file);
      const form = new FormData();
      for (const [key, value] of body.entries()) if (key !== 'file') form.append(key, value);
      form.append('file', compressed);
      return await super.upload(endpoint, form, orgId);
    } catch (error) {
      if (error && typeof error === 'object' && 'status' in error) throw error;
      throw { error: this.translate.instant('ACCOUNT.PHOTO_ERROR') };
    }
  }
}
