import { NextRequest, NextResponse } from 'next/server';
import { readFile } from 'fs/promises';
import { join } from 'path';

export async function GET(
  request: NextRequest,
  { params }: { params: { service: string } }
) {
  try {
    const { service } = params;
    const artifactsPath = join(process.cwd(), '.artifacts', 'openapi3', 'zitadel');
    
    // Try the .openapi.yaml format first
    let filePath: string;
    let content: string;
    
    try {
      filePath = join(artifactsPath, `${service}.openapi.yaml`);
      content = await readFile(filePath, 'utf-8');
      return new NextResponse(content, {
        headers: {
          'Content-Type': 'application/yaml',
          'Access-Control-Allow-Origin': '*',
        },
      });
    } catch {
      // Fallback to other possible formats
      try {
        filePath = join(artifactsPath, `${service}.yaml`);
        content = await readFile(filePath, 'utf-8');
        return new NextResponse(content, {
          headers: {
            'Content-Type': 'application/yaml',
            'Access-Control-Allow-Origin': '*',
          },
        });
      } catch {
        filePath = join(artifactsPath, `${service}.json`);
        content = await readFile(filePath, 'utf-8');
        return new NextResponse(content, {
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
          },
        });
      }
    }
  } catch (error) {
    console.error(`Error reading OpenAPI spec for service ${params.service}:`, error);
    return NextResponse.json(
      { error: `OpenAPI specification not found for service: ${params.service}` },
      { status: 404 }
    );
  }
}
