'use client';

import { useEffect, useState, useRef } from 'react';
import { createApiReference } from '@scalar/api-reference';
import yaml from 'js-yaml';

interface OpenApiSpec {
  name: string;
  fileName: string;
  content: string;
}

interface ApiResponse {
  specs: OpenApiSpec[];
}

export function ApiReferenceComponent() {
  const [specs, setSpecs] = useState<OpenApiSpec[]>([]);
  const [selectedSpec, setSelectedSpec] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    async function loadSpecs() {
      try {
        const response = await fetch('/api/openapi');
        if (!response.ok) {
          throw new Error('Failed to load API specifications');
        }
        const data: ApiResponse = await response.json();
        setSpecs(data.specs);
        if (data.specs.length > 0) {
          // Default to management service if available, otherwise first service
          const managementSpec = data.specs.find(spec => spec.name === 'management');
          setSelectedSpec(managementSpec ? 'management' : data.specs[0].name);
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    }

    loadSpecs();
  }, []);

  useEffect(() => {
    if (selectedSpec && containerRef.current) {
      const selectedSpecData = specs.find(spec => spec.name === selectedSpec);
      if (selectedSpecData) {
        try {
          const parsedSpec = yaml.load(selectedSpecData.content) as any;
          
          // Debug: Log the parsed spec
          console.log('Selected spec:', selectedSpec);
          console.log('Parsed spec:', parsedSpec);
          console.log('Parsed spec paths:', parsedSpec.paths);
          console.log('Number of paths:', Object.keys(parsedSpec.paths || {}).length);
          
          // Clear the container
          containerRef.current.innerHTML = '';
          
          // Create a div for Scalar
          const scalarDiv = document.createElement('div');
          scalarDiv.id = `api-reference-${selectedSpec}`;
          containerRef.current.appendChild(scalarDiv);
          
          // Create the API reference with correct configuration
          createApiReference(scalarDiv, {
            spec: {
              content: parsedSpec
            },
            theme: 'default',
          } as any);
        } catch (err) {
          console.error('Error parsing YAML or creating API reference:', err);
          if (containerRef.current) {
            containerRef.current.innerHTML = `<div style="padding: 20px; color: red;">Error loading API documentation: ${err}</div>`;
          }
        }
      }
    }
  }, [selectedSpec, specs]);

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontSize: '18px'
      }}>
        Loading API documentation...
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontSize: '18px',
        color: '#e74c3c'
      }}>
        <h2>Error loading API documentation</h2>
        <p>{error}</p>
        <p style={{ marginTop: '20px', fontSize: '14px', color: '#666' }}>
          Make sure to run <code>pnpm run generate</code> to generate the OpenAPI specifications.
        </p>
      </div>
    );
  }

  if (specs.length === 0) {
    return (
      <div style={{ 
        display: 'flex', 
        flexDirection: 'column',
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '100vh',
        fontSize: '18px'
      }}>
        <h2>No API specifications found</h2>
        <p style={{ marginTop: '20px', fontSize: '14px', color: '#666' }}>
          Run <code>pnpm run generate</code> to generate the OpenAPI specifications from proto files.
        </p>
      </div>
    );
  }

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      {/* Sidebar for service selection */}
      <div style={{ 
        width: '250px', 
        borderRight: '1px solid #e1e4e8',
        padding: '20px',
        overflowY: 'auto',
        backgroundColor: '#f6f8fa'
      }}>
        <h2 style={{ marginBottom: '20px', fontSize: '18px', fontWeight: 'bold' }}>
          ZITADEL API Services
        </h2>
        <nav>
          {specs.map((spec) => (
            <button
              key={spec.name}
              onClick={() => setSelectedSpec(spec.name)}
              style={{
                display: 'block',
                width: '100%',
                padding: '10px',
                marginBottom: '8px',
                border: 'none',
                borderRadius: '6px',
                backgroundColor: selectedSpec === spec.name ? '#0969da' : 'transparent',
                color: selectedSpec === spec.name ? 'white' : '#24292f',
                textAlign: 'left',
                cursor: 'pointer',
                fontSize: '14px',
                transition: 'background-color 0.2s'
              }}
              onMouseEnter={(e) => {
                if (selectedSpec !== spec.name) {
                  e.currentTarget.style.backgroundColor = '#f3f4f6';
                }
              }}
              onMouseLeave={(e) => {
                if (selectedSpec !== spec.name) {
                  e.currentTarget.style.backgroundColor = 'transparent';
                }
              }}
            >
              {spec.name.replace(/([A-Z])/g, ' $1').replace(/^./, (str) => str.toUpperCase())}
            </button>
          ))}
        </nav>
      </div>

      {/* Main content area */}
      <div 
        ref={containerRef}
        style={{ flex: 1, height: '100vh', overflow: 'auto' }}
      />
    </div>
  );
}
