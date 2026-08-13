import { ReactElement } from 'react';

export default function ApiUnderDevelopmentPage(): ReactElement {
    return (
        <div className="flex min-h-screen flex-col items-center justify-center p-8 text-center">
            <h1 className="mb-4 text-3xl font-bold">API Under Development</h1>
            <p className="text-muted-foreground max-w-md">
                The API endpoints for this resource are currently under development and are not ready to be used yet. Please check back later!
            </p>
        </div>
    );
}
