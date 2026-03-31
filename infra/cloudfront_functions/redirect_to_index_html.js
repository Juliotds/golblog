import cf from 'cloudfront';

const kvs = cf.kvs();

async function handler(event) {
    var request = event.request;
    var uri = request.uri;

    // Handle comment submission
    if (request.method === 'POST' && uri === '/comments') {
        var bodyText = request.body.encoding === 'base64'
            ? atob(request.body.data)
            : request.body.data;

        var params = new URLSearchParams(bodyText);
        var slug   = params.get('slug')   || '';
        var author = params.get('author') || 'anonymous';
        var body   = params.get('body')   || '';

        if (!slug || !body.trim()) {
            return {
                statusCode: 400,
                statusDescription: 'Bad Request',
                headers: { 'content-type': { value: 'text/plain' } },
                body: 'Missing slug or body.',
            };
        }

        var comments = [];
        try {
            var existing = await kvs.get(slug);
            if (existing) comments = JSON.parse(existing);
        } catch (_) {
            // Key not found — first comment for this slug
        }

        comments.push({
            author: author,
            body: body,
            date: new Date().toISOString().slice(0, 10),
        });

        await kvs.put(slug, JSON.stringify(comments));

        return {
            statusCode: 302,
            statusDescription: 'Found',
            headers: { location: { value: '/blog/' + slug + '#comments' } },
        };
    }

    // Default: append index.html to directory requests
    if (uri.endsWith('/')) {
        request.uri += 'index.html';
    } else if (!uri.includes('.')) {
        request.uri += '/index.html';
    }

    return request;
}
