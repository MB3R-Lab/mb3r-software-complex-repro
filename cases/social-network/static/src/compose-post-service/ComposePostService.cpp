#include <opentracing/tracer.h>

void ConfigureTracing() {
  SetUpTracer("src/compose-post-service/jaeger-config.yml", "compose-post-service");
  opentracing::Tracer::Global()->Inject(span->context(), writer);
}
