# Event Bus - Infrastructure

Este módulo contiene las implementaciones del event bus para diferentes backends.

## Implementaciones Disponibles

### In-Memory

Implementación en memoria para testing y desarrollo.

**Ubicación:** [`inmemory/`](./inmemory/README.md)

**Características:**
- Thread-safe
- Soporte para múltiples consumidores
- Delay mediante goroutines
- Sin persistencia

**Ideal para:**
- Testing unitario e integración
- Desarrollo local
- Prototipado rápido

### RabbitMQ

Implementación usando RabbitMQ para producción.

**Ubicación:** [`rabbit/`](./rabbit/README.md)

**Características:**
- Persistencia de mensajes
- Alta disponibilidad
- Delay usando DLX con TTL
- Broadcast con múltiples queues
- Soporte para microservicios

**Ideal para:**
- Producción
- Microservicios distribuidos
- Aplicaciones que requieren persistencia

## Elegir una Implementación

### Usa In-Memory si:
- ✅ Estás haciendo testing
- ✅ Desarrollando localmente sin infraestructura
- ✅ Tienes una aplicación simple y monolítica
- ✅ No necesitas persistencia

### Usa RabbitMQ si:
- ✅ Estás en producción
- ✅ Tienes múltiples servicios/microservicios
- ✅ Necesitas persistencia de eventos
- ✅ Necesitas alta disponibilidad
- ✅ Necesitas comunicación entre procesos/máquinas

## Ejemplos

### In-Memory

Ver `examples/infrastructure/eventbus/inmemory/main.go` para un ejemplo completo.

### RabbitMQ

Ver `examples/infrastructure/eventbus/rabbit/` para múltiples ejemplos:

- **Basic**: Publicación y consumo básico
- **Delayed**: Mensajes con delay
- **Signup**: Patrón broadcast con múltiples consumidores

## Referencias

- [Event Bus Core Documentation](../../core/eventbus/README.md)
- [In-Memory Implementation](./inmemory/README.md)
- [RabbitMQ Implementation](./rabbit/README.md)
