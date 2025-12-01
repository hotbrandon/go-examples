### Why check X-Request-ID even though RequestID() runs first?

Because clients may want to supply their own request IDs, especially:
✔ API clients that want to trace requests end-to-end

Mobile apps, frontends, or other microservices often generate a correlation ID, attach it, and then log it on their side too.

Example:

X-Request-ID: 3f2f7f16-d0f4-45c1-aed9-5266f90f1a91


Your backend receives it, reuses it, and includes it in logs.
This allows debugging across multiple services.

✔ Distributed tracing across services

If service A → service B → service C, you want a shared request ID for logs.

✔ Load balancers / API gateways may provide one

Cloud systems often inject their own:

Envoy: x-request-id

AWS ALB: x-amzn-trace-id

NGINX: X-Request-ID

Cloudflare: CF-RAY

You don’t want to overwrite these.

✔ Better developer experience

If a customer reports an issue and gives you a request ID, you want to be able to search for it through logs.
