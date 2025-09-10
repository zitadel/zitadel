import { NextRequest, NextResponse } from 'next/server';
import { readdir, readFile } from 'fs/promises';
import { join } from 'path';

export async function GET(request: NextRequest) {
  try {
    const artifactsPath = join(process.cwd(), '.artifacts', 'openapi3', 'zitadel');
    
    // Get all OpenAPI spec files
    const files = await readdir(artifactsPath);
    const openApiFiles = files.filter((file: string) => file.endsWith('.openapi.yaml'));
    
    const specs = await Promise.all(
      openApiFiles.map(async (file: string) => {
        const filePath = join(artifactsPath, file);
        const content = await readFile(filePath, 'utf-8');
        const serviceName = file.replace('.openapi.yaml', '');
        
        return {
          name: serviceName,
          fileName: file,
          content: content,
        };
      })
    );

    return NextResponse.json({ specs });
  } catch (error) {
    console.error('Error reading OpenAPI specs:', error);
    return NextResponse.json(
      { error: 'Failed to load OpenAPI specifications' },
      { status: 500 }
    );
  }
}
