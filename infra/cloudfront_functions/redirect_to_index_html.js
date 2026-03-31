import cf from 'cloudfront';

const kvs = cf.kvs('CommentsKVS');

const redirect = (path) => ({
    statusCode: 302,
    statusDescription: 'Found',
    headers: { location: { value: path } },
});

async function handler(event) {
    var request = event.request;
    var uri = request.uri;
    var method = request.method;

    // Block anything that isn't GET or POST /comments
    if (method !== 'GET' && !(method === 'POST' && uri === '/comments')) {
        return redirect('/invalid-operation');
    }

    // Handle comment submission
    if (method === 'POST' && uri === '/comments') {
        try {
            var bodyText = request.body.encoding === 'base64'
                ? atob(request.body.data)
                : request.body.data;

            var params = new URLSearchParams(bodyText);
            var slug   = params.get('slug')   || '';
            var author = params.get('author') || 'anonymous';
            var body   = params.get('body')   || '';

            if (!slug || !body.trim()) {
                return redirect('/invalid-operation');
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

            return redirect('/comment-posted');
        } catch (_) {
            return redirect('/error');
        }
    }

    // Default: append index.html to directory requests
    if (uri.endsWith('/')) {
        request.uri += 'index.html';
    } else if (!uri.includes('.')) {
        request.uri += '/index.html';
    }

    return request;
}
