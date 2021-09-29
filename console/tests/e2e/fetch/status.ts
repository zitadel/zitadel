class HTTPResponseError extends Error {
	constructor(response) {
		super(`HTTP Error Response: ${response.status} ${response.statusText}`);
	}
}

export function checkStatus(response) {
		// response.status >= 200 && response.status < 300
    if (!response.ok)
		throw new HTTPResponseError(response);
}
