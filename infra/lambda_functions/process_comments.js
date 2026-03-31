import { DynamoDBClient, GetItemCommand, PutItemCommand } from '@aws-sdk/client-dynamodb';

const client = new DynamoDBClient({ region: process.env.AWS_REGION || 'us-east-1' });
const TABLE  = process.env.COMMENTS_TABLE || 'Comments';

export const handler = async function(event) {
    const params = new URLSearchParams(event.rawQueryString || '');
    const slug   = params.get('slug')   || '';
    const author = params.get('author') || 'anonymous';
    const body   = params.get('body')   || '';

    if (!slug || !body.trim()) {
        return { statusCode: 302, headers: { Location: '/invalid-operation' } };
    }

    try {
        let comments = [];
        try {
            const { Item } = await client.send(new GetItemCommand({
                TableName: TABLE,
                Key: { slug: { S: slug } },
            }));
            if (Item && Item.comments) comments = JSON.parse(Item.comments.S);
        } catch (_) {
            // Item doesn't exist yet
        }

        comments.push({
            author,
            body,
            date: new Date().toISOString().slice(0, 10),
        });

        await client.send(new PutItemCommand({
            TableName: TABLE,
            Item: {
                slug:     { S: slug },
                comments: { S: JSON.stringify(comments) },
            },
        }));

        return { statusCode: 302, headers: { Location: '/comment-posted' } };
    } catch (err) {
        console.error('process_comments error:', err);
        return { statusCode: 302, headers: { Location: '/error' } };
    }
};
