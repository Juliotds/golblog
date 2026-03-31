const redirect = (path) => ({
    statusCode: 302,
    statusDescription: 'Found',
    headers: { location: { value: path } },
});

function qs(request, key) {
    var entry = request.querystring[key];
    return entry ? entry.value : '';
}

async function handler(event) {
    var request = event.request;
    var uri = request.uri;
    var method = request.method;

    // Block anything that isn't GET or POST /comments
    if (method !== 'GET' && !(method === 'POST' && uri === '/comments')) {
        return redirect('/invalid-operation');
    }

    // Pass POST /comments straight through to Lambda origin
    if (method === 'POST' && uri === '/comments') {
        return request;
    }

    // Default: append index.html to directory requests
    if (uri.endsWith('/')) {
        request.uri += 'index.html';
    } else if (!uri.includes('.')) {
        request.uri += '/index.html';
    }

    return request;
}
