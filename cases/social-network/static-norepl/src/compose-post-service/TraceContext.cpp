#include <opentracing/tracer.h>

void InjectTraceContext() {
  opentracing::Tracer::Global()->Inject(span->context(), writer);
}
