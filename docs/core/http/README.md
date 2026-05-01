# HTTP Client - Core

El módulo `core/http` proporciona abstracciones e interfaces para realizar peticiones HTTP de forma tipada y segura. Sigue los principios de Clean Architecture, separando las interfaces del core de las implementaciones concretas.

## 📋 Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Uso Básico](#uso-básico)
- [API](#api)
- [Tipos](#tipos)
- [Ejemplos](#ejemplos)
- [Casos de Uso](#casos-de-uso)

## ✨ Características

- **Interfaz abstracta**: `Client` interface que permite diferentes implementaciones
- **Respuestas tipadas**: Soporte para respuestas genéricas con `TypedResponse[T]`
- **Request/Response tipados**: Objetos estructurados para requests y responses
- **Serialización automática**: Marshal automático de objetos a JSON en requests
- **Funcional options**: Patrón de opciones funcionales para configuración
- **Criterios de éxito configurables**: Define qué códigos de estado se consideran exitosos
- **Medición de tiempo**: Tracking automático de duración de requests
- **Type-safe**: Uso de genéricos para seguridad de tipos

## 📦 Instalación

```bash
go get github.com/foundathyon/base/core/http
```

## 🚀 Uso Básico

### Request Simple

```go
import (
    "github.com/foundathyon/base/core/customctx"
    corehttp "github.com/foundathyon/base/core/http"
    "github.com/foundathyon/base/infrastructure/http"
)

ctx := customctx.New(context.Background())
client := http.NewClient(http.DefaultConfig())

// GET request simple
resp, err := client.Get(ctx, "https://api.example.com/users")
if err != nil {
    log.Fatal(err)
}
defer resp.Close()
```

### Request con Respuesta Tipada

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// GET con respuesta tipada
resp, err := corehttp.GetTyped[User](
    client,
    ctx,
    "https://api.example.com/users/1",
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %+v\n", resp.Body)
fmt.Printf("Status: %d\n", resp.StatusCode)
fmt.Printf("Duration: %v\n", resp.Duration)
```

### POST con Body Tipado

```go
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type UserResponse struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

createReq := CreateUserRequest{
    Name:  "John Doe",
    Email: "john@example.com",
}

// POST con request y response tipados
// TRequest = CreateUserRequest, TResponse = UserResponse
resp, err := corehttp.PostTyped[CreateUserRequest, UserResponse](
    client,
    ctx,
    "https://api.example.com/users",
    createReq, // El objeto se serializa automáticamente a JSON
    corehttp.WithContentType("application/json"),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created User: %+v\n", resp.Body)
```

## 📚 API

### Client Interface

La interfaz `Client` define los métodos para realizar peticiones HTTP:

```go
type Client interface {
    Do(ctx *customctx.CustomContext, req *Request) (*Response, error)
    Get(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
    Post(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)
    Put(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)
    Patch(ctx *customctx.CustomContext, url string, body io.Reader, opts ...RequestOption) (*Response, error)
    Delete(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
    Head(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
    Options(ctx *customctx.CustomContext, url string, opts ...RequestOption) (*Response, error)
}
```

### Request

Estructura que representa una petición HTTP:

```go
type Request struct {
    Method      string
    URL         string
    Headers     map[string]string
    Body        io.Reader
    ContentType string
    QueryParams map[string]string
    Timeout     *int // en segundos
}
```

### Response

Estructura que representa una respuesta HTTP:

```go
type Response struct {
    StatusCode    int
    Status        string
    Headers       map[string]string
    Body          io.ReadCloser
    ContentLength int64
}
```

### TypedResponse

Wrapper genérico para respuestas tipadas:

```go
type TypedResponse[T any] struct {
    StatusCode              int
    Status                  string
    Headers                 map[string]string
    Body                    T              // Respuesta parseada del tipo T
    RawBody                 []byte         // Body crudo (útil para binarios)
    ContentType             string
    ContentLength           int64
    RequestTime             time.Time
    ResponseTime            time.Time
    Duration                time.Duration
    ExpectedStatusCode      *int
    SuccessStatusCodeRange  *StatusCodeRange
}
```

### Funciones Helper Tipadas

Funciones genéricas para realizar peticiones con respuestas tipadas:

```go
// GET con respuesta tipada
func GetTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)

// POST con request y response tipados
func PostTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error)

// PUT con request y response tipados
func PutTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error)

// PATCH con request y response tipados
func PatchTyped[TRequest any, TResponse any](client Client, ctx *customctx.CustomContext, url string, body TRequest, opts ...RequestOption) (*TypedResponse[TResponse], error)

// DELETE con respuesta tipada
func DeleteTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)

// HEAD con respuesta tipada
func HeadTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)

// OPTIONS con respuesta tipada
func OptionsTyped[T any](client Client, ctx *customctx.CustomContext, url string, opts ...RequestOption) (*TypedResponse[T], error)
```

### Request Options

Opciones funcionales para configurar requests:

```go
// Añadir un header
WithHeader(key, value string) RequestOption

// Añadir múltiples headers
WithHeaders(headers map[string]string) RequestOption

// Añadir un query parameter
WithQueryParam(key, value string) RequestOption

// Añadir múltiples query parameters
WithQueryParams(params map[string]string) RequestOption

// Establecer Content-Type
WithContentType(contentType string) RequestOption

// Establecer timeout en segundos
WithTimeout(seconds int) RequestOption
```

### TypedResponse Options

Opciones para configurar criterios de éxito:

```go
// Establecer código de estado esperado
WithExpectedStatusCode[T](code int) func(*TypedResponse[T])

// Establecer rango de códigos de estado exitosos
WithSuccessStatusCodeRange[T](min, max int) func(*TypedResponse[T])
```

### RequestBody Interface

Interfaz para tipos que pueden serializarse a request body:

```go
type RequestBody interface {
    ToReader() (io.Reader, error)
}
```

Si un tipo implementa `RequestBody`, se usará su método `ToReader()`. Si no, se intentará hacer JSON marshal automático.

## 📖 Tipos

### RequestBody

Interfaz para serialización personalizada de request bodies:

```go
type RequestBody interface {
    ToReader() (io.Reader, error)
}
```

**Implementaciones incluidas:**
- `NewJSONRequestBody(value any) RequestBody` - Serializa cualquier valor a JSON
- `NewReaderRequestBody(reader io.Reader) RequestBody` - Envuelve un `io.Reader` existente

### StatusCodeRange

Rango de códigos de estado HTTP:

```go
type StatusCodeRange struct {
    Min int // Mínimo (inclusivo)
    Max int // Máximo (inclusivo)
}
```

## 💡 Ejemplos

### Ejemplo 1: GET Simple

```go
ctx := customctx.New(context.Background())
client := http.NewClient(http.DefaultConfig())

resp, err := client.Get(ctx, "https://api.example.com/data")
if err != nil {
    log.Fatal(err)
}
defer resp.Close()

body, _ := resp.ReadBodyString()
fmt.Println(body)
```

### Ejemplo 2: GET con Respuesta Tipada

```go
type Product struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

resp, err := corehttp.GetTyped[Product](
    client,
    ctx,
    "https://api.example.com/products/1",
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Product: %+v\n", resp.Body)
fmt.Printf("Is Success: %v\n", resp.IsSuccess())
```

### Ejemplo 3: POST con Body Automático

```go
type CreateOrder struct {
    UserID    int    `json:"user_id"`
    ProductID int    `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

type Order struct {
    ID       int    `json:"id"`
    UserID   int    `json:"user_id"`
    Total    float64 `json:"total"`
    Status   string  `json:"status"`
}

order := CreateOrder{
    UserID:    123,
    ProductID: 456,
    Quantity:  2,
}

resp, err := corehttp.PostTyped[CreateOrder, Order](
    client,
    ctx,
    "https://api.example.com/orders",
    order, // Se serializa automáticamente a JSON
    corehttp.WithContentType("application/json"),
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Order created: %+v\n", resp.Body)
```

### Ejemplo 4: Criterios de Éxito Personalizados

```go
// Esperar específicamente código 201
resp, err := corehttp.PostTyped[CreateUser, User](
    client,
    ctx,
    "https://api.example.com/users",
    user,
    corehttp.WithExpectedStatusCode[User](201),
)
if err != nil {
    log.Fatal(err)
}

if resp.IsSuccess() {
    fmt.Println("User created successfully!")
}

// O usar un rango
resp.SetSuccessStatusCodeRange(&corehttp.StatusCodeRange{Min: 200, Max: 299})
```

### Ejemplo 5: Request con Query Parameters

```go
resp, err := corehttp.GetTyped[[]Product](
    client,
    ctx,
    "https://api.example.com/products",
    corehttp.WithQueryParam("category", "electronics"),
    corehttp.WithQueryParam("limit", "10"),
    corehttp.WithQueryParam("offset", "0"),
)
```

### Ejemplo 6: Request con Timeout

```go
resp, err := corehttp.GetTyped[Data](
    client,
    ctx,
    "https://api.example.com/data",
    corehttp.WithTimeout(5), // 5 segundos
)
```

### Ejemplo 7: Respuesta Binaria

```go
// Para respuestas binarias, usa []byte como tipo
resp, err := corehttp.GetTyped[[]byte](
    client,
    ctx,
    "https://api.example.com/image.jpg",
)
if err != nil {
    log.Fatal(err)
}

// El body ya está en resp.Body como []byte
ioutil.WriteFile("image.jpg", resp.Body, 0644)
```

## 🎯 Casos de Uso

### Integración con APIs REST

El módulo está diseñado para facilitar la integración con APIs REST externas:

```go
type APIClient struct {
    httpClient corehttp.Client
    baseURL    string
}

func (c *APIClient) GetUser(id int) (*User, error) {
    ctx := customctx.New(context.Background())
    resp, err := corehttp.GetTyped[User](
        c.httpClient,
        ctx,
        fmt.Sprintf("%s/users/%d", c.baseURL, id),
    )
    if err != nil {
        return nil, err
    }
    return &resp.Body, nil
}
```

### Testing con Mocks

La interfaz `Client` permite crear mocks fácilmente para testing:

```go
type MockClient struct{}

func (m *MockClient) Get(ctx *customctx.CustomContext, url string, opts ...corehttp.RequestOption) (*corehttp.Response, error) {
    // Implementación mock
}
```

### Medición de Performance

El `TypedResponse` incluye información de timing:

```go
resp, err := corehttp.GetTyped[Data](client, ctx, url)
if err != nil {
    return err
}

fmt.Printf("Request took: %v\n", resp.Duration)
fmt.Printf("Started at: %v\n", resp.RequestTime)
fmt.Printf("Completed at: %v\n", resp.ResponseTime)
```

## 🔗 Referencias

- [Infrastructure HTTP Implementation](../infrastructure/http/README.md) - Implementación concreta usando `net/http`
- [Ejemplos](../../../examples/infrastructure/http/) - Ejemplos de uso completos
