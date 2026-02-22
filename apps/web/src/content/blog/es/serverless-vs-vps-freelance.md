---
title: "Serverless vs. VPS: lo que nadie te dice cuando trabajas solo"
description: "Después de años eligiendo infraestructura para proyectos reales, comparto mi proceso para decidir entre serverless y un VPS — y por qué terminé donde terminé."
pubDate: 2026-02-15
tags: ["serverless", "vps", "arquitectura", "freelance", "devops", "go", "svelte"]
category: "opinion"
lang: "es"
draft: false
---

Hace un par de años me senté a rediseñar la forma en la que entregaba proyectos. Llevaba tiempo trabajando con hosting compartido, algún que otro VPS configurado a medias, y ese sentimiento constante de que "debería estar usando algo más moderno". Serverless estaba en todas las conversaciones: Lambda, Cloud Functions, contenedores gestionados. La narrativa era clara — si no estás en la nube, te estás quedando atrás.

Así que hice lo que cualquier desarrollador curioso haría: me metí de lleno. Configuré proyectos en AWS, experimenté con Lambda y API Gateway, probé DynamoDB, armé pipelines con Step Functions. Y aprendí mucho. Pero también aprendí algo que nadie mencionaba en los tutoriales ni en los threads de Twitter: **que serverless no está diseñado para resolver los problemas que yo tenía**.

Este artículo no es una comparación genérica de "pros y contras". Es lo que descubrí después de ir y venir entre ambos mundos, y por qué mi infraestructura actual se ve como se ve.

## El problema de partida

Cuando trabajas solo o en un equipo muy reducido, tu recurso más escaso no es el cómputo ni el almacenamiento — es tu tiempo. Cada hora que dedicas a infraestructura es una hora que no estás construyendo features, atendiendo clientes, o simplemente descansando para no quemarte.

Y aquí es donde empecé a notar una desconexión entre lo que se recomienda en la industria y lo que necesitaba en la práctica. La mayoría de los artículos sobre arquitectura cloud están escritos desde la perspectiva de equipos con al menos un DevOps dedicado, presupuestos de infraestructura de cuatro cifras mensuales, y un volumen de tráfico que justifica la complejidad. Pero cuando tu cliente es una empresa local con 50 empleados que necesita un sistema interno, ese contexto simplemente no aplica.

No estoy diciendo que serverless sea malo. Estoy diciendo que la pregunta correcta no es "¿serverless o VPS?" sino "¿qué problema estoy resolviendo y cuál es la forma más directa de resolverlo?"

## Lo que descubrí al usar serverless en proyectos reales

Mi primer proyecto serio con Lambda fue un sistema de notificaciones. La lógica era sencilla: recibir un evento, procesar datos, enviar un email o un mensaje. En papel, el caso de uso perfecto para serverless — funciones cortas, event-driven, sin estado.

Y funcionó. Pero el camino para llegar ahí fue revelador.

Configurar los permisos de IAM me tomó más tiempo del que esperaba. No porque sea difícil conceptualmente, sino porque el modelo de permisos de AWS es absurdamente granular. Cada función necesitaba un rol, cada rol necesitaba políticas específicas, y un error en cualquier política resultaba en un mensaje críptico que te mandaba a buscar en Stack Overflow qué permiso faltaba. Para un equipo con un ingeniero de seguridad que se encarga de eso, no es problema. Para mí, solo, era fricción pura.

Después vino el debugging. Cuando algo fallaba en Lambda, el flujo era: revisar CloudWatch logs (que tienen su propio pricing), encontrar el request ID, buscar el error, hacer un cambio, deployar, esperar el cold start, y probar de nuevo. Compara eso con tener un servidor donde haces `docker logs` y ves todo en tiempo real, o metes un breakpoint y depuras directo.

El costo también fue sorpresa. El cómputo de Lambda es barato, sí. Pero API Gateway cobra por request. Si necesitas una base de datos relacional, RDS tiene un costo base mensual que no baja de los $15 USD aunque no la uses. Si tu Lambda necesita acceder a recursos dentro de una VPC, necesitas un NAT Gateway — que cuesta alrededor de $30 USD/mes solo por existir. CloudWatch te cobra por ingesta de logs. De pronto, un proyecto que en un VPS de $10 corría completo, en AWS me estaba costando $50-70 mensuales. Y eso sin contar mi tiempo configurando todo.

La [documentación de pricing de AWS](https://aws.amazon.com/lambda/pricing/) es transparente sobre los costos de cómputo, pero los costos colaterales — el Gateway, el NAT, los logs, el almacenamiento de artefactos — son los que te sorprenden.

## El momento de la reflexión

Después de tres o cuatro proyectos así, tuve que sentarme a pensar con honestidad. ¿Estaba usando serverless porque resolvía un problema real, o porque quería tenerlo en mi currículum? La respuesta fue incómoda.

La realidad es que la mayoría de mis clientes tienen tráfico predecible. No hay picos de Black Friday. No hay millones de requests concurrentes. Hay 30, 50, tal vez 200 usuarios accediendo a un sistema interno durante horas de oficina. Para ese patrón de uso, un servidor que está encendido 24/7 por $15-25 dólares al mes no es desperdicio — es simplicidad.

Y la simplicidad tiene un valor enorme que es difícil de cuantificar. Cuando algo falla a las 11 de la noche (porque siempre falla a las 11 de la noche), quiero poder conectarme por SSH, ver los logs, identificar el problema, y resolverlo. No quiero abrir la consola de AWS, navegar entre 6 servicios diferentes, descubrir que el error está en un paso de una Step Function que triggerea una Lambda que escribe en DynamoDB. Quiero ver un stack trace claro en un solo lugar.

## Entonces, ¿serverless no sirve?

No es eso. Serverless resuelve problemas reales y concretos. Los resuelve muy bien. Pero esos problemas son específicos:

**Cargas de trabajo que escalan de cero a miles y vuelven a cero.** Si tienes un evento que genera tráfico masivo por unas horas y después nada, pagar por un servidor 24/7 sería desperdiciar dinero. Serverless escala automáticamente y pagas solo lo que usas. Esto es particularmente valioso para aplicaciones como procesamiento de archivos en batch, donde recibes un lote grande de documentos una vez al día y el resto del tiempo no hay nada que procesar.

**Funciones aisladas que no justifican un servidor propio.** Un webhook que recibe notificaciones de un servicio externo y las procesa. Un cron que genera un reporte PDF cada lunes. Un endpoint que redimensiona imágenes al subirlas. Estas son funciones que se ejecutan esporádicamente, duran segundos, y no necesitan mantener estado. Es el sweet spot de Lambda y Cloud Functions.

**Integraciones dentro de un ecosistema cloud existente.** Si tu cliente ya tiene toda su infraestructura en AWS y quiere agregar una funcionalidad nueva, pelear contra el ecosistema no tiene sentido. Usas lo que ya está ahí.

El problema es cuando tomas estas ventajas y las extrapolas a toda tu arquitectura. He visto (y cometido el error de armar) aplicaciones enteras donde cada endpoint es una Lambda, el estado se maneja entre DynamoDB y S3, y la comunicación entre servicios pasa por SQS y SNS. Funciona, pero la complejidad operacional es desproporcionada para lo que el proyecto necesita.

## El camino de regreso al VPS (con lo aprendido)

Cuando decidí volver a centrar mi infraestructura en VPS, no fue un retroceso. Fue aplicar todo lo que había aprendido de serverless — deploys automatizados, infraestructura reproducible, servicios desacoplados — pero sobre una base que podía manejar solo.

El catalizador fue Docker. Todo lo que me gustaba del modelo serverless — empaquetar mi código con sus dependencias, deployar de forma reproducible, escalar servicios de forma independiente — lo podía hacer con contenedores en un VPS. Sin la capa de abstracción de un cloud provider, sin el pricing complejo, y con control total sobre mi entorno.

Mi setup actual es un VPS con Docker Compose, Traefik como reverse proxy con SSL automático vía Let's Encrypt, y cada servicio corriendo en su propio contenedor. El deploy es un `docker compose up -d --build` después de un push a main. GitHub Actions se encarga del CI. No hay magia, no hay vendor lock-in, y si mañana quiero mover todo a otro proveedor de VPS, es copiar archivos y correr el mismo comando.

Traefik merece mención especial porque resolvió algo que con Nginx me tomaba configuración manual: el ruteo automático basado en labels de Docker y la renovación de certificados SSL. Definir que `api.midominio.com` apunte al servicio correcto es una label en el `docker-compose.yml`, no un archivo de configuración aparte. La [documentación de Traefik](https://doc.traefik.io/traefik/) es de las mejores que he visto en herramientas de infraestructura.

## Eligiendo el stack: por qué Go y Svelte

Esta parte del proceso fue más difícil de lo que esperaba, porque involucra matar vacas sagradas.

Venía de PHP y TypeScript. Conocía Laravel, conocía Express, conocía Next.js. ¿Por qué cambiar algo que ya funcionaba? La razón fue puramente práctica: quería correr múltiples microservicios en un VPS sin que se comieran toda la RAM.

Un servicio de Laravel con PHP-FPM consume entre 80 y 150 MB de RAM en reposo. Uno de Express/Node entre 60 y 120 MB. No es un problema si tienes un monolito. Pero si tu arquitectura tiene 4 o 5 servicios independientes, ya estás en 400-600 MB solo en runtime, sin contar la base de datos. En un VPS de 2 GB, eso es apretado.

Go cambió esa ecuación. Un servicio en Go compilado consume entre 10 y 25 MB de RAM sirviendo tráfico real. Puedo correr 5 o 6 microservicios, una base de datos, un reverse proxy y un servicio de monitoreo en un VPS de 2 GB y todavía me sobra más de un gig de RAM. Eso para un freelance significa poder hostear varios proyectos de clientes en una misma máquina si es necesario, o tener margen enorme para crecer.

Pero la eficiencia de recursos no fue la única razón. Go tiene algo que no encontré en otros lenguajes de backend: una legibilidad a prueba del tiempo. Cuando abro un archivo de Go que escribí hace 8 meses, entiendo qué hace en minutos. No hay decorators que investigar, no hay inyección de dependencias mágica, no hay middleware implícito. El flujo del programa es lineal y explícito. Recibes un request, validas, procesas, respondes. El manejo de errores con `if err != nil` parece verboso al principio, pero después de un tiempo te das cuenta de que es exactamente la claridad que quieres cuando estás debuggeando a las 11 de la noche.

El router que uso es [Chi](https://github.com/go-chi/chi), que sigue la interfaz estándar de `net/http` de Go. Esto significa que no estoy aprendiendo un framework propietario — estoy usando la librería estándar con un poco de conveniencia encima. Si Chi dejara de mantenerse mañana, migrar a otro router o incluso al `net/http` puro sería cambiar unas cuantas líneas.

Para el frontend, la decisión fue más personal. Había usado React por años y la fatiga era real. No la fatiga del framework en sí, sino la del ecosistema que cambia cada seis meses. Next.js cambiando su modelo de rendering, la transición a server components, el debate de App Router vs Pages Router, los problemas de hidratación... demasiado movimiento para algo que debería ser estable.

Svelte fue un respiro. Un componente de Svelte es HTML con reactividad declarativa. No hay virtual DOM, no hay hooks con reglas crípticas, no hay `useEffect` con arrays de dependencias que nunca sabes si están bien. Los estilos son scoped al componente por defecto — escribes CSS normal y no se filtra a otros componentes. Y con Svelte 5, el modelo de reactividad con runes (`$state`, `$derived`, `$effect`) es tan directo que el código casi se lee como pseudocódigo.

Combinado con [Astro](https://astro.build) como framework de sitios, tengo lo mejor de ambos mundos: páginas estáticas donde no necesito interactividad (blog, landing, páginas informativas) y componentes Svelte hidratados solo donde hacen falta (formularios, dashboards, elementos interactivos). El resultado es un sitio que carga casi instantáneo porque el HTML ya viene renderizado del servidor, y JavaScript solo se carga donde es necesario.

## La base de datos: SQLite y cuándo sí escalar a PostgreSQL

Otra decisión que va contra la corriente: uso SQLite para la mayoría de mis proyectos.

Sé que suena raro. SQLite tiene fama de ser "la base de datos de desarrollo" o "la que usas para prototipar". Pero esa reputación está desactualizada. SQLite en modo WAL (Write-Ahead Logging) maneja lecturas concurrentes sin problema. Para aplicaciones con un patrón de lectura intensiva y escritura moderada — que es exactamente lo que la mayoría de los sistemas internos son — rinde igual o mejor que PostgreSQL, con la ventaja de que no hay un servidor de base de datos que configurar, mantener, o que pueda caerse.

La base de datos es un archivo. Lo respaldas copiándolo (o con herramientas como [Litestream](https://litestream.io) que hacen replicación continua a S3). Lo migras copiándolo. No hay connection pooling que configurar, no hay usuarios de base de datos que gestionar, no hay puertos que exponer.

¿Cuándo escalo a PostgreSQL? Cuando necesito escrituras concurrentes pesadas desde múltiples procesos, búsqueda full-text avanzada, tipos de datos especializados como JSON indexado o geoespacial, o cuando el volumen de datos supera lo razonable para un solo archivo. Para la gran mayoría de proyectos de empresas pequeñas y medianas, ese punto no llega nunca.

## El deploy: GitHub Actions, Docker, y nada más

Mi pipeline de deploy es deliberadamente simple. Push a `main`, GitHub Actions construye las imágenes Docker, las sube al registro, y el VPS hace pull y reinicia los servicios. No hay Kubernetes, no hay Terraform, no hay Ansible. Un `docker-compose.yml` define toda la infraestructura.

¿Es esto "lo correcto" según las mejores prácticas de DevOps? Probablemente no, si estuviéramos hablando de una empresa con 20 servicios y un equipo de plataforma. Pero para un freelance con 3-6 servicios por proyecto, es exactamente lo que necesito: algo que entiendo completamente, que puedo debuggear sin documentación, y que no requiere mantener infraestructura auxiliar.

El Dockerfile de un servicio Go típico es mínimo:

```dockerfile
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /server .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

La imagen final pesa entre 15 y 25 MB dependiendo de las dependencias. Compara eso con una imagen de Node que fácilmente supera los 200 MB, o una de PHP-FPM que ronda los 150 MB. La diferencia se nota en los tiempos de deploy y en cuántas imágenes puedes tener en el registro sin que el almacenamiento se dispare.

## Los conflictos que sigo teniendo

No quiero pintar esto como si fuera la solución perfecta, porque no lo es. Hay cosas que me siguen generando fricción.

**Go no tiene un ORM que me convenza del todo.** GORM es popular pero agrega una capa de magia que va contra la filosofía de Go. SQLx es más explícito pero requiere escribir SQL a mano para todo. Terminé escribiendo SQL directo con la librería estándar `database/sql`, y aunque funciona bien, a veces extraño la productividad de Eloquent en Laravel o Prisma en TypeScript. Es un tradeoff consciente: más control y transparencia a cambio de más líneas de código en la capa de datos.

**Svelte tiene un ecosistema más chico que React.** Cuando necesitas un componente de terceros — un date picker complejo, un editor de texto rico, una librería de gráficos — las opciones en Svelte son menos y a veces menos maduras. Esto es cada vez menos problema conforme el ecosistema crece, pero es una realidad que hay que considerar.

**El monitoreo en un VPS es manual.** En AWS tienes CloudWatch, X-Ray, dashboards automáticos. En un VPS, si quieres monitoreo, lo configuras tú. Personalmente uso structured logging con `slog` (la librería estándar de Go desde la versión 1.21) y lo mando a stdout para que Docker lo capture. No es tan sofisticado como un APM completo, pero para mis necesidades es suficiente. Herramientas como [Dozzle](https://dozzle.dev/) te dan un dashboard de logs basado en web sin configuración, lo cual ayuda bastante.

**La escalabilidad horizontal tiene un techo.** Si un proyecto crece al punto de necesitar múltiples servidores, balanceo de carga, y failover automático, un solo VPS con Docker Compose se queda corto. Pero en mi experiencia, ese punto llega mucho más tarde de lo que la mayoría cree. Un VPS de 4 GB bien optimizado con Go sirve cómodamente miles de requests por segundo. Cuando un proyecto llega a ese nivel de tráfico, generalmente ya tiene presupuesto para una infraestructura más robusta.

## El estado actual

Hoy mi infraestructura estándar para un proyecto nuevo se ve así: un VPS en Hetzner o DigitalOcean ($15-25 USD/mes), Docker Compose para orquestación, Traefik para el reverse proxy y SSL, servicios de backend en Go con Chi, frontend en Astro con componentes Svelte, SQLite para persistencia con Litestream para backups, y GitHub Actions para CI/CD.

El total de RAM que usa un proyecto típico con 2-3 microservicios, el frontend, la base de datos y Traefik ronda los 200-300 MB. Eso me deja margen enorme en un VPS de 2 GB, y la posibilidad de hostear más de un proyecto por máquina si los clientes son pequeños.

¿Uso serverless? Sí, pero para lo que tiene sentido. Un webhook que recibe datos de un servicio externo y los procesa. Un cron pesado que no vale la pena tener corriendo en un contenedor. Funciones puntuales donde el modelo de pay-per-invocation realmente ahorra. Pero no es la base de mi arquitectura — es un complemento.

## Lo que le diría a alguien que está decidiendo

Si estás empezando como freelance o estás reconsiderando tu infraestructura, mi consejo es que resistas la presión de adoptar lo que está de moda y te enfoques en lo que te permite entregar rápido y mantener con confianza.

No necesitas Kubernetes para servir una aplicación con 200 usuarios. No necesitas una arquitectura de microservicios distribuidos en Lambda para un CRUD con reportes. No necesitas DynamoDB para una tabla de 10,000 registros. Pero sí necesitas entender qué pasa cuando algo se rompe, y poder arreglarlo sin consultar tres documentaciones diferentes.

Elige la herramienta que entiendas de cabo a rabo. Si eso es PHP con Laravel en un hosting compartido, funciona. Si es Node en un VPS, funciona. Si es Go con Docker, funciona. La tecnología importa menos de lo que la industria quiere hacerte creer. Lo que importa es que puedas entregar, que puedas mantener, y que tu cliente esté contento.

Y si después de leer todo esto sigues queriendo probar serverless — hazlo. Pero hazlo en una función aislada de un proyecto real, no en todo un sistema. Así aprendes el modelo sin apostar todo tu proyecto a una arquitectura que tal vez no necesitas.

Al final del día, la mejor infraestructura es la que te deja enfocarte en el código y en los problemas de tu cliente, no en la infraestructura misma.
