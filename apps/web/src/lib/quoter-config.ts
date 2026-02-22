type Labels = { es: string; en: string };

export interface PaymentPlan {
  key: string;
  label: Labels;
  description: Labels;
  icon: string;
  badge?: Labels;
  calculate: (total: number, currency: string) => GeneratedPlan;
}

export interface GeneratedPlan {
  key: string;
  label: Labels;
  description: Labels;
  icon: string;
  badge?: Labels;
  primary: string;
  secondary?: string;
  details: Labels;
  totalCost: number;
}

export const SOURCE_CODE_SURCHARGE = 0.25;

export interface ProjectType {
  key: string;
  base: number;
  label: Labels;
  icon: string;
  description?: Labels;
  image?: string;
  features: string[];
}

export interface Feature {
  key: string;
  cost: number;
  label: Labels;
  description?: Labels;
  icon?: string;
}

export const PROJECT_TYPES: ProjectType[] = [
  {
    key: 'website',
    base: 7500,
    label: { es: 'Página Web', en: 'Website' },
    icon: 'globe',
    description: {
      es: 'Tu presencia digital profesional. Incluye diseño personalizado, optimización para móviles, SEO, blog y formularios de contacto. Ideal para captar clientes en línea.',
      en: 'Your professional digital presence. Includes custom design, mobile optimization, SEO, blog and contact forms. Perfect for attracting clients online.',
    },
    image: 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'responsiveDesign',
      'blog',
      'contactForm',
      'seo',
      'multiLang',
      'adminPanel',
      'socialMedia',
      'animations',
      'analytics',
      'liveChat',
    ],
  },
  {
    key: 'ecommerce',
    base: 17500,
    label: { es: 'Tienda en Línea', en: 'Online Store' },
    icon: 'shoppingCart',
    description: {
      es: 'Vende productos en línea con carrito de compras, pagos con tarjeta, integración de envíos, facturación automática y gestión de inventario.',
      en: 'Sell products online with shopping cart, card payments, shipping integration, automatic invoicing and inventory management.',
    },
    image: 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'productCatalog',
      'shoppingCart',
      'stripePayments',
      'paypalPayments',
      'shippingIntegration',
      'cfdiEcommerce',
      'inventorySync',
      'couponsDiscounts',
      'productReviews',
      'wishlist',
      'orderTracking',
      'emailMarketing',
    ],
  },
  {
    key: 'mobileApp',
    base: 20000,
    label: { es: 'Aplicación Móvil', en: 'Mobile App' },
    icon: 'smartphone',
    description: {
      es: 'App nativa para iOS y Android con Flutter o React Native. Push notifications, modo offline, GPS, biometría y publicación en tiendas de apps.',
      en: 'Native app for iOS and Android with Flutter or React Native. Push notifications, offline mode, GPS, biometrics and app store publishing.',
    },
    image: 'https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'crossPlatform',
      'pushNotifications',
      'offlineMode',
      'gpsLocation',
      'biometricAuth',
      'cameraIntegration',
      'appStorePublish',
      'inAppPayments',
      'deepLinking',
      'socialLogin',
    ],
  },
  {
    key: 'adminSystem',
    base: 17500,
    label: { es: 'Sistema Administrativo', en: 'Management System' },
    icon: 'monitor',
    description: {
      es: 'Sistema a medida para gestionar tu negocio: usuarios, reportes, exportación de datos, multi-sucursal, auditoría y flujos de trabajo automatizados.',
      en: 'Custom system to manage your business: users, reports, data export, multi-branch, auditing and automated workflows.',
    },
    image: 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'usersRoles',
      'reports',
      'exportExcelPdf',
      'emailNotifications',
      'multiBranch',
      'auditLog',
      'externalApi',
      'docGeneration',
      'workflows',
      'payroll',
    ],
  },
  {
    key: 'saas',
    base: 25000,
    label: { es: 'Plataforma SaaS', en: 'SaaS Platform' },
    icon: 'cloud',
    description: {
      es: 'Plataforma multi-cliente con suscripciones, API pública, webhooks, métricas de uso y marca blanca. Lista para escalar tu negocio digital.',
      en: 'Multi-tenant platform with subscriptions, public API, webhooks, usage metrics and white labeling. Ready to scale your digital business.',
    },
    image: 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'multiTenant',
      'subscriptions',
      'onboarding',
      'publicApi',
      'webhooks',
      'customDomains',
      'usageMetrics',
      'teamManagement',
      'rolePermissions',
      'whiteLabeling',
    ],
  },
  {
    key: 'inventory',
    base: 12500,
    label: { es: 'Control de Inventario', en: 'Inventory Management' },
    icon: 'package',
    description: {
      es: 'Controla entradas, salidas, alertas de stock bajo, códigos de barra, multi-almacén y órdenes de compra en tiempo real.',
      en: 'Track entries, exits, low stock alerts, barcodes, multi-warehouse and purchase orders in real time.',
    },
    image: 'https://images.unsplash.com/photo-1553413077-190dd305871c?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'stockInOut',
      'lowStockAlerts',
      'barcodeQr',
      'movementReports',
      'multiWarehouse',
      'posIntegration',
      'batchTracking',
      'purchaseOrders',
    ],
  },
  {
    key: 'pos',
    base: 15000,
    label: { es: 'Punto de Venta', en: 'Point of Sale' },
    icon: 'dollarSign',
    description: {
      es: 'Punto de venta con registro de ventas, corte de caja, múltiples métodos de pago, tickets, descuentos y reportes de ventas.',
      en: 'Point of sale with sales registry, cash cuts, multiple payment methods, tickets, discounts and sales reports.',
    },
    image: 'https://images.unsplash.com/photo-1556742111-a301076d9d18?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'salesRegistry',
      'cashCut',
      'multiPayment',
      'tickets',
      'discounts',
      'salesReports',
      'loyaltyProgram',
      'vendorControl',
    ],
  },
  {
    key: 'billing',
    base: 14000,
    label: { es: 'Facturación Automática', en: 'Automated Billing' },
    icon: 'fileText',
    description: {
      es: 'Genera CFDI, facturación recurrente, portal de clientes, reportes fiscales e integración contable. Cumple con el SAT al 100%.',
      en: 'Generate CFDI invoices, recurring billing, client portal, tax reports and accounting integration. Fully SAT compliant.',
    },
    image: 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'cfdiGeneration',
      'satCatalog',
      'recurringBilling',
      'clientPortal',
      'taxReports',
      'accountingIntegration',
      'massBilling',
      'creditNotes',
    ],
  },
  {
    key: 'booking',
    base: 15000,
    label: { es: 'Reservaciones / Citas', en: 'Bookings / Appointments' },
    icon: 'calendarCheck',
    description: {
      es: 'Sistema de reservaciones en línea con calendario, recordatorios SMS, agenda por empleado y sincronización con Google Calendar.',
      en: 'Online booking system with calendar view, SMS reminders, employee schedules and Google Calendar sync.',
    },
    image: 'https://images.unsplash.com/photo-1506784983877-45594efa4cbe?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'onlineBooking',
      'calendarView',
      'smsReminders',
      'employeeSchedule',
      'googleCalendarSync',
      'waitlist',
      'recurringBookings',
      'depositPayments',
    ],
  },
  {
    key: 'apiIntegration',
    base: 10000,
    label: { es: 'API e Integraciones', en: 'API & Integrations' },
    icon: 'link',
    description: {
      es: 'Conecta tus sistemas con APIs REST/GraphQL, Stripe, Twilio, PayPal, SAT, Uber Direct, Envia.com y OAuth/SSO.',
      en: 'Connect your systems with REST/GraphQL APIs, Stripe, Twilio, PayPal, SAT, Uber Direct, Envia.com and OAuth/SSO.',
    },
    image: 'https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'restApi',
      'stripeIntegration',
      'twilioIntegration',
      'paypalIntegration',
      'satCfdiApi',
      'uberDirectApi',
      'enviacomApi',
      'oauthSso',
      'graphqlApi',
      'apiDocs',
    ],
  },
  {
    key: 'cloudDevOps',
    base: 12500,
    label: { es: 'Infraestructura Cloud / DevOps', en: 'Cloud Infrastructure / DevOps' },
    icon: 'server',
    description: {
      es: 'Infraestructura en la nube con AWS, Docker, CI/CD, Terraform, SSL, monitoreo, auto-scaling y respaldos automáticos.',
      en: 'Cloud infrastructure with AWS, Docker, CI/CD, Terraform, SSL, monitoring, auto-scaling and automatic backups.',
    },
    image: 'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'awsSetup',
      'dockerContainers',
      'ciCdPipeline',
      'terraformIac',
      'sslCerts',
      'monitoringAlerts',
      'autoScaling',
      'lambdaFunctions',
      'backupStrategy',
      'loadBalancing',
    ],
  },
  {
    key: 'techUpdate',
    base: 6000,
    label: { es: 'Actualización Tecnológica', en: 'Tech Modernization' },
    icon: 'refreshCw',
    description: {
      es: 'Moderniza tu infraestructura: diagnóstico, configuración de equipos, migración de datos, capacitación y soporte técnico continuo.',
      en: 'Modernize your infrastructure: diagnostics, device setup, data migration, staff training and ongoing technical support.',
    },
    image: 'https://images.unsplash.com/photo-1518770660439-4636190af475?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'infraDiagnostic',
      'deviceSetup',
      'networkConfig',
      'dataMigration',
      'staffTraining',
      'postSupport',
      'cloudMigration',
      'securityAudit',
    ],
  },
  {
    key: 'aiIntegration',
    base: 10000,
    label: { es: 'Integración con IA', en: 'AI Integration' },
    icon: 'bot',
    description: {
      es: 'Integra inteligencia artificial: chatbots para WhatsApp y web, automatización de procesos, análisis de datos, asistentes virtuales y procesamiento de documentos.',
      en: 'Integrate artificial intelligence: WhatsApp and web chatbots, process automation, data analysis, virtual assistants and document processing.',
    },
    image: 'https://images.unsplash.com/photo-1677442136019-21780ecad995?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'whatsappBot',
      'webChatbot',
      'processAutomation',
      'dataAnalysis',
      'smartReports',
      'virtualAssistant',
      'voiceAssistant',
      'docProcessing',
    ],
  },
  {
    key: 'consulting',
    base: 10000,
    label: { es: 'Consultoría y Arquitectura', en: 'Consulting & Architecture' },
    icon: 'compass',
    description: {
      es: 'Análisis técnico, diseño de arquitectura, revisión de código y roadmap tecnológico para tu proyecto.',
      en: 'Technical analysis, architecture design, code review and technology roadmap for your project.',
    },
    image: 'https://images.unsplash.com/photo-1552664730-d307ca884978?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'techAudit',
      'archDesign',
      'codeReview',
      'roadmap',
      'stackSelection',
      'perfOptimization',
      'scalabilityPlan',
      'docAndDiagrams',
    ],
  },
  {
    key: 'teamTraining',
    base: 8000,
    label: { es: 'Capacitación de Equipos', en: 'Team Training' },
    icon: 'graduationCap',
    description: {
      es: 'Entrena a tu equipo en nuevas tecnologías, herramientas y procesos. Cursos presenciales o remotos, adaptados a tu stack.',
      en: 'Train your team on new technologies, tools and processes. On-site or remote courses, adapted to your stack.',
    },
    image: 'https://images.unsplash.com/photo-1524178232363-1fb2b075b655?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'needsAssessment',
      'customCurriculum',
      'liveWorkshops',
      'trainingMaterials',
      'practiceProjects',
      'postTrainingSupport',
      'certificationPath',
      'recordedSessions',
    ],
  },
  {
    key: 'migration',
    base: 15000,
    label: { es: 'Migración y Modernización de Software', en: 'Software Migration & Modernization' },
    icon: 'arrowUpRight',
    description: {
      es: 'Lleva tu sistema legacy a tecnologías modernas: refactorización de código, migración a la nube y optimización de rendimiento.',
      en: 'Bring your legacy system to modern technologies: code refactoring, cloud migration and performance optimization.',
    },
    image: 'https://images.unsplash.com/photo-1504868584819-f8e8b4b6d7e3?w=600&h=320&fit=crop&auto=format&q=80',
    features: [
      'legacyAudit',
      'codeRefactor',
      'dbMigration',
      'cloudMigrationMod',
      'apiModernization',
      'testingSetup',
      'perfTuning',
      'documentationMod',
    ],
  },
];

export const FEATURES: Record<string, Feature> = {
  // ── Website ──────────────────────────────────────────────
  responsiveDesign: {
    key: 'responsiveDesign', cost: 0, icon: 'monitor',
    label: { es: 'Diseño responsive', en: 'Responsive design' },
    description: { es: 'Se adapta a móvil, tablet y escritorio', en: 'Adapts to mobile, tablet and desktop' },
  },
  blog: {
    key: 'blog', cost: 2500, icon: 'fileText',
    label: { es: 'Blog integrado', en: 'Integrated blog' },
    description: { es: 'Publica artículos y noticias en tu sitio', en: 'Publish articles and news on your site' },
  },
  contactForm: {
    key: 'contactForm', cost: 1000, icon: 'mail',
    label: { es: 'Formulario de contacto', en: 'Contact form' },
    description: { es: 'Recibe mensajes directamente de tu sitio', en: 'Receive messages directly from your site' },
  },
  seo: {
    key: 'seo', cost: 2500, icon: 'search',
    label: { es: 'SEO optimizado', en: 'SEO optimized' },
    description: { es: 'Mejor posicionamiento en Google y buscadores', en: 'Better ranking on Google and search engines' },
  },
  multiLang: {
    key: 'multiLang', cost: 4000, icon: 'globe',
    label: { es: 'Multi-idioma', en: 'Multi-language' },
    description: { es: 'Contenido en español, inglés u otros idiomas', en: 'Content in Spanish, English or other languages' },
  },
  adminPanel: {
    key: 'adminPanel', cost: 6000, icon: 'settings',
    label: { es: 'Panel de administración', en: 'Admin panel' },
    description: { es: 'Gestiona contenido sin saber programar', en: 'Manage content without coding knowledge' },
  },
  socialMedia: {
    key: 'socialMedia', cost: 1500, icon: 'share',
    label: { es: 'Integración con redes sociales', en: 'Social media integration' },
    description: { es: 'Conecta con Facebook, Instagram, X y más', en: 'Connect with Facebook, Instagram, X and more' },
  },
  animations: {
    key: 'animations', cost: 2000, icon: 'sparkles',
    label: { es: 'Animaciones y efectos visuales', en: 'Animations & visual effects' },
    description: { es: 'Transiciones suaves y micro-interacciones', en: 'Smooth transitions and micro-interactions' },
  },
  analytics: {
    key: 'analytics', cost: 1500, icon: 'barChart',
    label: { es: 'Google Analytics / métricas', en: 'Google Analytics / metrics' },
    description: { es: 'Mide visitas, conversiones y comportamiento', en: 'Track visits, conversions and behavior' },
  },
  liveChat: {
    key: 'liveChat', cost: 2500, icon: 'messageCircle',
    label: { es: 'Chat en vivo', en: 'Live chat' },
    description: { es: 'Atiende a tus visitantes en tiempo real', en: 'Assist your visitors in real time' },
  },

  // ── E-commerce ───────────────────────────────────────────
  productCatalog: {
    key: 'productCatalog', cost: 0, icon: 'grid',
    label: { es: 'Catálogo de productos', en: 'Product catalog' },
    description: { es: 'Organiza y muestra tus productos con filtros', en: 'Organize and display products with filters' },
  },
  shoppingCart: {
    key: 'shoppingCart', cost: 0, icon: 'shoppingCart',
    label: { es: 'Carrito de compras', en: 'Shopping cart' },
    description: { es: 'Carrito persistente con resumen de orden', en: 'Persistent cart with order summary' },
  },
  stripePayments: {
    key: 'stripePayments', cost: 3000, icon: 'creditCard',
    label: { es: 'Pagos con Stripe', en: 'Stripe payments' },
    description: { es: 'Tarjetas de crédito/débito internacionales', en: 'International credit/debit cards' },
  },
  paypalPayments: {
    key: 'paypalPayments', cost: 2500, icon: 'creditCard',
    label: { es: 'Pagos con PayPal', en: 'PayPal payments' },
    description: { es: 'Pagos seguros vía PayPal', en: 'Secure payments via PayPal' },
  },
  shippingIntegration: {
    key: 'shippingIntegration', cost: 4000, icon: 'truck',
    label: { es: 'Integración de envíos', en: 'Shipping integration' },
    description: { es: 'Envia.com, Uber Direct, tarifas automáticas', en: 'Envia.com, Uber Direct, automatic rates' },
  },
  cfdiEcommerce: {
    key: 'cfdiEcommerce', cost: 3500, icon: 'fileText',
    label: { es: 'Facturación CFDI automática', en: 'Automatic CFDI invoicing' },
    description: { es: 'Genera facturas al momento de la compra', en: 'Generate invoices at purchase time' },
  },
  inventorySync: {
    key: 'inventorySync', cost: 3000, icon: 'refreshCw',
    label: { es: 'Sincronización de inventario', en: 'Inventory sync' },
    description: { es: 'Stock actualizado en tiempo real', en: 'Stock updated in real time' },
  },
  couponsDiscounts: {
    key: 'couponsDiscounts', cost: 2000, icon: 'tag',
    label: { es: 'Cupones y descuentos', en: 'Coupons & discounts' },
    description: { es: 'Códigos promocionales y ofertas especiales', en: 'Promo codes and special offers' },
  },
  productReviews: {
    key: 'productReviews', cost: 1500, icon: 'star',
    label: { es: 'Reseñas de productos', en: 'Product reviews' },
    description: { es: 'Los clientes califican y comentan productos', en: 'Customers rate and review products' },
  },
  wishlist: {
    key: 'wishlist', cost: 1500, icon: 'heart',
    label: { es: 'Lista de deseos', en: 'Wishlist' },
    description: { es: 'Guarda productos favoritos para después', en: 'Save favorite products for later' },
  },
  orderTracking: {
    key: 'orderTracking', cost: 2500, icon: 'mapPin',
    label: { es: 'Rastreo de pedidos', en: 'Order tracking' },
    description: { es: 'Seguimiento en tiempo real del envío', en: 'Real-time shipment tracking' },
  },
  emailMarketing: {
    key: 'emailMarketing', cost: 2500, icon: 'mail',
    label: { es: 'Email marketing automatizado', en: 'Automated email marketing' },
    description: { es: 'Carritos abandonados, newsletters, promos', en: 'Abandoned carts, newsletters, promos' },
  },

  // ── Mobile App ───────────────────────────────────────────
  crossPlatform: {
    key: 'crossPlatform', cost: 0, icon: 'smartphone',
    label: { es: 'Multiplataforma (iOS + Android)', en: 'Cross-platform (iOS + Android)' },
    description: { es: 'Flutter o React Native, una sola base de código', en: 'Flutter or React Native, single codebase' },
  },
  pushNotifications: {
    key: 'pushNotifications', cost: 2500, icon: 'bell',
    label: { es: 'Notificaciones push', en: 'Push notifications' },
    description: { es: 'Envía alertas y recordatorios al dispositivo', en: 'Send alerts and reminders to device' },
  },
  offlineMode: {
    key: 'offlineMode', cost: 4000, icon: 'wifi',
    label: { es: 'Modo offline', en: 'Offline mode' },
    description: { es: 'Funciona sin conexión, sincroniza después', en: 'Works offline, syncs later' },
  },
  gpsLocation: {
    key: 'gpsLocation', cost: 3000, icon: 'mapPin',
    label: { es: 'GPS / geolocalización', en: 'GPS / geolocation' },
    description: { es: 'Mapas, rutas y ubicación en tiempo real', en: 'Maps, routes and real-time location' },
  },
  biometricAuth: {
    key: 'biometricAuth', cost: 2500, icon: 'shield',
    label: { es: 'Autenticación biométrica', en: 'Biometric authentication' },
    description: { es: 'Huella digital, Face ID, reconocimiento facial', en: 'Fingerprint, Face ID, facial recognition' },
  },
  cameraIntegration: {
    key: 'cameraIntegration', cost: 2000, icon: 'camera',
    label: { es: 'Integración de cámara', en: 'Camera integration' },
    description: { es: 'Fotos, escaneo de documentos o QR', en: 'Photos, document or QR scanning' },
  },
  appStorePublish: {
    key: 'appStorePublish', cost: 4000, icon: 'upload',
    label: { es: 'Publicación en tiendas', en: 'App store publishing' },
    description: { es: 'Google Play Store y Apple App Store', en: 'Google Play Store and Apple App Store' },
  },
  inAppPayments: {
    key: 'inAppPayments', cost: 3500, icon: 'creditCard',
    label: { es: 'Pagos in-app', en: 'In-app payments' },
    description: { es: 'Compras y suscripciones dentro de la app', en: 'Purchases and subscriptions within the app' },
  },
  deepLinking: {
    key: 'deepLinking', cost: 1500, icon: 'link',
    label: { es: 'Deep linking', en: 'Deep linking' },
    description: { es: 'URLs que abren secciones específicas de la app', en: 'URLs that open specific app sections' },
  },
  socialLogin: {
    key: 'socialLogin', cost: 2000, icon: 'users',
    label: { es: 'Login con redes sociales', en: 'Social login' },
    description: { es: 'Inicia sesión con Google, Apple, Facebook', en: 'Sign in with Google, Apple, Facebook' },
  },

  // ── Admin System ─────────────────────────────────────────
  usersRoles: {
    key: 'usersRoles', cost: 4000, icon: 'shield',
    label: { es: 'Control de usuarios y roles', en: 'User and role management' },
    description: { es: 'Permisos granulares por rol y usuario', en: 'Granular permissions by role and user' },
  },
  reports: {
    key: 'reports', cost: 3500, icon: 'barChart',
    label: { es: 'Reportes y dashboards', en: 'Reports and dashboards' },
    description: { es: 'Visualiza métricas clave de tu negocio', en: 'Visualize key business metrics' },
  },
  exportExcelPdf: {
    key: 'exportExcelPdf', cost: 2000, icon: 'download',
    label: { es: 'Exportar a Excel/PDF', en: 'Export to Excel/PDF' },
    description: { es: 'Descarga reportes en formatos estándar', en: 'Download reports in standard formats' },
  },
  emailNotifications: {
    key: 'emailNotifications', cost: 1500, icon: 'bell',
    label: { es: 'Notificaciones por email', en: 'Email notifications' },
    description: { es: 'Alertas automáticas por correo electrónico', en: 'Automatic email alerts' },
  },
  multiBranch: {
    key: 'multiBranch', cost: 5000, icon: 'building',
    label: { es: 'Multi-sucursal', en: 'Multi-branch' },
    description: { es: 'Gestiona varias ubicaciones desde un sistema', en: 'Manage multiple locations from one system' },
  },
  auditLog: {
    key: 'auditLog', cost: 2500, icon: 'clipboard',
    label: { es: 'Auditoría / logs de actividad', en: 'Audit / activity logs' },
    description: { es: 'Historial de quién hizo qué y cuándo', en: 'History of who did what and when' },
  },
  externalApi: {
    key: 'externalApi', cost: 3000, icon: 'link',
    label: { es: 'API para integraciones externas', en: 'API for external integrations' },
    description: { es: 'Conecta con otros sistemas y servicios', en: 'Connect with other systems and services' },
  },
  docGeneration: {
    key: 'docGeneration', cost: 2500, icon: 'fileText',
    label: { es: 'Generación de documentos', en: 'Document generation' },
    description: { es: 'Contratos, cotizaciones y reportes en PDF', en: 'Contracts, quotes and reports in PDF' },
  },
  workflows: {
    key: 'workflows', cost: 4000, icon: 'gitBranch',
    label: { es: 'Flujos de trabajo automatizados', en: 'Automated workflows' },
    description: { es: 'Aprobaciones, escalaciones y tareas automáticas', en: 'Approvals, escalations and automatic tasks' },
  },
  payroll: {
    key: 'payroll', cost: 5000, icon: 'dollarSign',
    label: { es: 'Nómina', en: 'Payroll' },
    description: { es: 'Cálculo de sueldos, deducciones e IMSS', en: 'Salary calculation, deductions and benefits' },
  },

  // ── SaaS ─────────────────────────────────────────────────
  multiTenant: {
    key: 'multiTenant', cost: 0, icon: 'layers',
    label: { es: 'Arquitectura multi-tenant', en: 'Multi-tenant architecture' },
    description: { es: 'Cada cliente con datos aislados y seguros', en: 'Each client with isolated, secure data' },
  },
  subscriptions: {
    key: 'subscriptions', cost: 4000, icon: 'creditCard',
    label: { es: 'Suscripciones y planes', en: 'Subscriptions & plans' },
    description: { es: 'Cobros recurrentes con Stripe o PayPal', en: 'Recurring billing with Stripe or PayPal' },
  },
  onboarding: {
    key: 'onboarding', cost: 3000, icon: 'userPlus',
    label: { es: 'Onboarding de usuarios', en: 'User onboarding' },
    description: { es: 'Guía paso a paso para nuevos usuarios', en: 'Step-by-step guide for new users' },
  },
  publicApi: {
    key: 'publicApi', cost: 3500, icon: 'code',
    label: { es: 'API pública', en: 'Public API' },
    description: { es: 'Permite a terceros integrarse con tu plataforma', en: 'Allow third parties to integrate with your platform' },
  },
  webhooks: {
    key: 'webhooks', cost: 2500, icon: 'link',
    label: { es: 'Webhooks', en: 'Webhooks' },
    description: { es: 'Notificaciones en tiempo real a otros sistemas', en: 'Real-time notifications to other systems' },
  },
  customDomains: {
    key: 'customDomains', cost: 4000, icon: 'globe',
    label: { es: 'Dominios personalizados', en: 'Custom domains' },
    description: { es: 'Cada cliente con su propio dominio', en: 'Each client with their own domain' },
  },
  usageMetrics: {
    key: 'usageMetrics', cost: 3000, icon: 'barChart',
    label: { es: 'Métricas de uso', en: 'Usage metrics' },
    description: { es: 'Dashboard de uso, límites y consumo', en: 'Usage dashboard, limits and consumption' },
  },
  teamManagement: {
    key: 'teamManagement', cost: 2500, icon: 'users',
    label: { es: 'Gestión de equipos', en: 'Team management' },
    description: { es: 'Invita miembros, asigna roles por equipo', en: 'Invite members, assign roles by team' },
  },
  rolePermissions: {
    key: 'rolePermissions', cost: 2500, icon: 'lock',
    label: { es: 'Roles y permisos granulares', en: 'Granular roles & permissions' },
    description: { es: 'Control fino de acceso por funcionalidad', en: 'Fine-grained access control by feature' },
  },
  whiteLabeling: {
    key: 'whiteLabeling', cost: 5000, icon: 'palette',
    label: { es: 'Marca blanca', en: 'White labeling' },
    description: { es: 'Personaliza colores, logo y dominio por cliente', en: 'Customize colors, logo and domain per client' },
  },

  // ── Inventory ────────────────────────────────────────────
  stockInOut: {
    key: 'stockInOut', cost: 0, icon: 'package',
    label: { es: 'Entradas y salidas de stock', en: 'Stock entries and exits' },
    description: { es: 'Registra movimientos de mercancía', en: 'Record merchandise movements' },
  },
  lowStockAlerts: {
    key: 'lowStockAlerts', cost: 1500, icon: 'bell',
    label: { es: 'Alertas de stock bajo', en: 'Low stock alerts' },
    description: { es: 'Notificaciones cuando el inventario es bajo', en: 'Notifications when inventory is low' },
  },
  barcodeQr: {
    key: 'barcodeQr', cost: 2500, icon: 'scan',
    label: { es: 'Códigos de barra / QR', en: 'Barcode / QR codes' },
    description: { es: 'Escaneo rápido para entradas y salidas', en: 'Fast scanning for entries and exits' },
  },
  movementReports: {
    key: 'movementReports', cost: 2000, icon: 'barChart',
    label: { es: 'Reportes de movimiento', en: 'Movement reports' },
    description: { es: 'Historial detallado de movimientos de stock', en: 'Detailed stock movement history' },
  },
  multiWarehouse: {
    key: 'multiWarehouse', cost: 4000, icon: 'building',
    label: { es: 'Multi-almacén', en: 'Multi-warehouse' },
    description: { es: 'Gestiona inventario en múltiples ubicaciones', en: 'Manage inventory across multiple locations' },
  },
  posIntegration: {
    key: 'posIntegration', cost: 3000, icon: 'dollarSign',
    label: { es: 'Integración con punto de venta', en: 'POS integration' },
    description: { es: 'Sincroniza ventas con tu inventario', en: 'Sync sales with your inventory' },
  },
  batchTracking: {
    key: 'batchTracking', cost: 2500, icon: 'layers',
    label: { es: 'Rastreo por lotes / caducidad', en: 'Batch tracking / expiration' },
    description: { es: 'Controla lotes, fechas de caducidad y FIFO', en: 'Track batches, expiration dates and FIFO' },
  },
  purchaseOrders: {
    key: 'purchaseOrders', cost: 3000, icon: 'clipboard',
    label: { es: 'Órdenes de compra', en: 'Purchase orders' },
    description: { es: 'Genera y rastrea pedidos a proveedores', en: 'Generate and track supplier orders' },
  },

  // ── POS ──────────────────────────────────────────────────
  salesRegistry: {
    key: 'salesRegistry', cost: 0, icon: 'dollarSign',
    label: { es: 'Registro de ventas', en: 'Sales registry' },
    description: { es: 'Punto de venta rápido y fácil de usar', en: 'Fast and easy-to-use point of sale' },
  },
  cashCut: {
    key: 'cashCut', cost: 1500, icon: 'clipboard',
    label: { es: 'Corte de caja', en: 'Cash cut' },
    description: { es: 'Cuadra efectivo al inicio y cierre del día', en: 'Balance cash at day open and close' },
  },
  multiPayment: {
    key: 'multiPayment', cost: 2500, icon: 'creditCard',
    label: { es: 'Múltiples métodos de pago', en: 'Multiple payment methods' },
    description: { es: 'Efectivo, tarjeta, transferencia, vales', en: 'Cash, card, transfer, vouchers' },
  },
  tickets: {
    key: 'tickets', cost: 1500, icon: 'fileText',
    label: { es: 'Tickets / recibos', en: 'Tickets / receipts' },
    description: { es: 'Impresión de tickets personalizados', en: 'Custom ticket printing' },
  },
  discounts: {
    key: 'discounts', cost: 1000, icon: 'tag',
    label: { es: 'Descuentos y promociones', en: 'Discounts and promotions' },
    description: { es: 'Aplica descuentos por producto o venta', en: 'Apply discounts by product or sale' },
  },
  salesReports: {
    key: 'salesReports', cost: 2000, icon: 'barChart',
    label: { es: 'Reportes de ventas', en: 'Sales reports' },
    description: { es: 'Métricas diarias, semanales y mensuales', en: 'Daily, weekly and monthly metrics' },
  },
  loyaltyProgram: {
    key: 'loyaltyProgram', cost: 3000, icon: 'award',
    label: { es: 'Programa de lealtad', en: 'Loyalty program' },
    description: { es: 'Puntos, recompensas y clientes frecuentes', en: 'Points, rewards and frequent customers' },
  },
  vendorControl: {
    key: 'vendorControl', cost: 2500, icon: 'users',
    label: { es: 'Control de vendedores', en: 'Vendor control' },
    description: { es: 'Comisiones, metas y desempeño por vendedor', en: 'Commissions, goals and performance per vendor' },
  },

  // ── Billing ──────────────────────────────────────────────
  cfdiGeneration: {
    key: 'cfdiGeneration', cost: 0, icon: 'fileText',
    label: { es: 'Generación de CFDI', en: 'CFDI generation' },
    description: { es: 'Facturas electrónicas válidas ante el SAT', en: 'Electronic invoices valid before SAT' },
  },
  satCatalog: {
    key: 'satCatalog', cost: 2000, icon: 'database',
    label: { es: 'Catálogo de productos SAT', en: 'SAT product catalog' },
    description: { es: 'Claves de producto y unidad del SAT', en: 'SAT product and unit codes' },
  },
  recurringBilling: {
    key: 'recurringBilling', cost: 3000, icon: 'refreshCw',
    label: { es: 'Facturación recurrente', en: 'Recurring billing' },
    description: { es: 'Genera facturas periódicas automáticamente', en: 'Generate periodic invoices automatically' },
  },
  clientPortal: {
    key: 'clientPortal', cost: 4000, icon: 'monitor',
    label: { es: 'Portal de descarga para clientes', en: 'Client download portal' },
    description: { es: 'Tus clientes descargan sus facturas en línea', en: 'Your clients download their invoices online' },
  },
  taxReports: {
    key: 'taxReports', cost: 2500, icon: 'barChart',
    label: { es: 'Reportes fiscales', en: 'Tax reports' },
    description: { es: 'Resumen de impuestos para contabilidad', en: 'Tax summary for accounting' },
  },
  accountingIntegration: {
    key: 'accountingIntegration', cost: 3500, icon: 'calculator',
    label: { es: 'Integración contable', en: 'Accounting integration' },
    description: { es: 'Conecta con CONTPAQi, Aspel u otros', en: 'Connect with accounting software' },
  },
  massBilling: {
    key: 'massBilling', cost: 2500, icon: 'files',
    label: { es: 'Facturación masiva', en: 'Mass billing' },
    description: { es: 'Genera cientos de facturas de una vez', en: 'Generate hundreds of invoices at once' },
  },
  creditNotes: {
    key: 'creditNotes', cost: 2000, icon: 'fileText',
    label: { es: 'Notas de crédito', en: 'Credit notes' },
    description: { es: 'Cancela o ajusta facturas emitidas', en: 'Cancel or adjust issued invoices' },
  },

  // ── Booking ──────────────────────────────────────────────
  onlineBooking: {
    key: 'onlineBooking', cost: 0, icon: 'calendarCheck',
    label: { es: 'Reservaciones en línea', en: 'Online bookings' },
    description: { es: 'Tus clientes reservan desde tu sitio web', en: 'Your clients book from your website' },
  },
  calendarView: {
    key: 'calendarView', cost: 2500, icon: 'calendar',
    label: { es: 'Vista de calendario', en: 'Calendar view' },
    description: { es: 'Visualiza todas las citas en un calendario', en: 'View all appointments in a calendar' },
  },
  smsReminders: {
    key: 'smsReminders', cost: 2500, icon: 'messageCircle',
    label: { es: 'Recordatorios SMS', en: 'SMS reminders' },
    description: { es: 'Envía recordatorios automáticos por SMS', en: 'Send automatic reminders via SMS' },
  },
  employeeSchedule: {
    key: 'employeeSchedule', cost: 3000, icon: 'users',
    label: { es: 'Agenda por empleado', en: 'Employee schedule' },
    description: { es: 'Cada empleado con su propia agenda', en: 'Each employee with their own schedule' },
  },
  googleCalendarSync: {
    key: 'googleCalendarSync', cost: 2000, icon: 'refreshCw',
    label: { es: 'Sincronización Google Calendar', en: 'Google Calendar sync' },
    description: { es: 'Sincroniza citas con Google Calendar', en: 'Sync appointments with Google Calendar' },
  },
  waitlist: {
    key: 'waitlist', cost: 1500, icon: 'clock',
    label: { es: 'Lista de espera', en: 'Waitlist' },
    description: { es: 'Gestiona clientes en espera automáticamente', en: 'Manage waiting clients automatically' },
  },
  recurringBookings: {
    key: 'recurringBookings', cost: 2500, icon: 'refreshCw',
    label: { es: 'Reservaciones recurrentes', en: 'Recurring bookings' },
    description: { es: 'Citas semanales, quincenales o mensuales', en: 'Weekly, biweekly or monthly appointments' },
  },
  depositPayments: {
    key: 'depositPayments', cost: 3000, icon: 'creditCard',
    label: { es: 'Pagos de anticipo', en: 'Deposit payments' },
    description: { es: 'Cobra anticipos al momento de reservar', en: 'Collect deposits at booking time' },
  },

  // ── API & Integrations ───────────────────────────────────
  restApi: {
    key: 'restApi', cost: 0, icon: 'code',
    label: { es: 'API REST', en: 'REST API' },
    description: { es: 'Endpoints seguros y documentados', en: 'Secure and documented endpoints' },
  },
  stripeIntegration: {
    key: 'stripeIntegration', cost: 3000, icon: 'creditCard',
    label: { es: 'Integración Stripe', en: 'Stripe integration' },
    description: { es: 'Pagos, suscripciones y payouts con Stripe', en: 'Payments, subscriptions and payouts with Stripe' },
  },
  twilioIntegration: {
    key: 'twilioIntegration', cost: 2500, icon: 'messageCircle',
    label: { es: 'Integración Twilio', en: 'Twilio integration' },
    description: { es: 'SMS, WhatsApp y llamadas automatizadas', en: 'SMS, WhatsApp and automated calls' },
  },
  paypalIntegration: {
    key: 'paypalIntegration', cost: 2500, icon: 'creditCard',
    label: { es: 'Integración PayPal', en: 'PayPal integration' },
    description: { es: 'Pagos y checkout con PayPal', en: 'Payments and checkout with PayPal' },
  },
  satCfdiApi: {
    key: 'satCfdiApi', cost: 3500, icon: 'fileText',
    label: { es: 'API SAT / CFDI', en: 'SAT / CFDI API' },
    description: { es: 'Timbrado y validación de facturas CFDI', en: 'CFDI invoice stamping and validation' },
  },
  uberDirectApi: {
    key: 'uberDirectApi', cost: 3000, icon: 'truck',
    label: { es: 'Uber Direct', en: 'Uber Direct' },
    description: { es: 'Entregas locales vía Uber Direct', en: 'Local deliveries via Uber Direct' },
  },
  enviacomApi: {
    key: 'enviacomApi', cost: 2500, icon: 'truck',
    label: { es: 'Envia.com', en: 'Envia.com' },
    description: { es: 'Envíos nacionales con múltiples paqueterías', en: 'National shipping with multiple carriers' },
  },
  oauthSso: {
    key: 'oauthSso', cost: 3000, icon: 'lock',
    label: { es: 'OAuth / SSO', en: 'OAuth / SSO' },
    description: { es: 'Single Sign-On con Google, Microsoft, etc.', en: 'Single Sign-On with Google, Microsoft, etc.' },
  },
  graphqlApi: {
    key: 'graphqlApi', cost: 2500, icon: 'code',
    label: { es: 'GraphQL API', en: 'GraphQL API' },
    description: { es: 'API flexible con queries optimizadas', en: 'Flexible API with optimized queries' },
  },
  apiDocs: {
    key: 'apiDocs', cost: 1500, icon: 'fileText',
    label: { es: 'Documentación de API', en: 'API documentation' },
    description: { es: 'Docs interactivos estilo Swagger / OpenAPI', en: 'Interactive docs Swagger / OpenAPI style' },
  },

  // ── Cloud / DevOps ───────────────────────────────────────
  awsSetup: {
    key: 'awsSetup', cost: 0, icon: 'cloud',
    label: { es: 'Configuración AWS', en: 'AWS setup' },
    description: { es: 'EC2, S3, RDS y servicios de Amazon', en: 'EC2, S3, RDS and Amazon services' },
  },
  dockerContainers: {
    key: 'dockerContainers', cost: 2500, icon: 'box',
    label: { es: 'Docker / contenedores', en: 'Docker / containers' },
    description: { es: 'Entornos reproducibles y portátiles', en: 'Reproducible and portable environments' },
  },
  ciCdPipeline: {
    key: 'ciCdPipeline', cost: 3000, icon: 'gitBranch',
    label: { es: 'CI/CD pipeline', en: 'CI/CD pipeline' },
    description: { es: 'Despliegue automático con cada commit', en: 'Automatic deployment with each commit' },
  },
  terraformIac: {
    key: 'terraformIac', cost: 3500, icon: 'terminal',
    label: { es: 'Terraform / IaC', en: 'Terraform / IaC' },
    description: { es: 'Infraestructura como código, versionada', en: 'Infrastructure as code, versioned' },
  },
  sslCerts: {
    key: 'sslCerts', cost: 1000, icon: 'shield',
    label: { es: 'Certificados SSL', en: 'SSL certificates' },
    description: { es: 'HTTPS automático con Let\'s Encrypt', en: 'Automatic HTTPS with Let\'s Encrypt' },
  },
  monitoringAlerts: {
    key: 'monitoringAlerts', cost: 2500, icon: 'activity',
    label: { es: 'Monitoreo y alertas', en: 'Monitoring & alerts' },
    description: { es: 'Uptime, métricas y alertas en tiempo real', en: 'Uptime, metrics and real-time alerts' },
  },
  autoScaling: {
    key: 'autoScaling', cost: 4000, icon: 'trending',
    label: { es: 'Auto-scaling', en: 'Auto-scaling' },
    description: { es: 'Escala recursos automáticamente con demanda', en: 'Scale resources automatically with demand' },
  },
  lambdaFunctions: {
    key: 'lambdaFunctions', cost: 3000, icon: 'zap',
    label: { es: 'Funciones serverless', en: 'Serverless functions' },
    description: { es: 'AWS Lambda, sin administrar servidores', en: 'AWS Lambda, no server management' },
  },
  backupStrategy: {
    key: 'backupStrategy', cost: 2000, icon: 'hardDrive',
    label: { es: 'Estrategia de respaldos', en: 'Backup strategy' },
    description: { es: 'Respaldos automáticos diarios en la nube', en: 'Automatic daily cloud backups' },
  },
  loadBalancing: {
    key: 'loadBalancing', cost: 2500, icon: 'server',
    label: { es: 'Balanceo de carga', en: 'Load balancing' },
    description: { es: 'Distribuye tráfico entre múltiples servidores', en: 'Distribute traffic across multiple servers' },
  },

  // ── Tech Update ──────────────────────────────────────────
  infraDiagnostic: {
    key: 'infraDiagnostic', cost: 0, icon: 'search',
    label: { es: 'Diagnóstico de infraestructura', en: 'Infrastructure diagnostic' },
    description: { es: 'Evaluación completa de tu infraestructura actual', en: 'Complete assessment of your current infrastructure' },
  },
  deviceSetup: {
    key: 'deviceSetup', cost: 1500, icon: 'monitor',
    label: { es: 'Configuración de equipos', en: 'Device setup' },
    description: { es: 'Instalación y configuración de hardware', en: 'Hardware installation and configuration' },
  },
  networkConfig: {
    key: 'networkConfig', cost: 2000, icon: 'wifi',
    label: { es: 'Configuración de red', en: 'Network configuration' },
    description: { es: 'WiFi, VPN, firewalls y seguridad de red', en: 'WiFi, VPN, firewalls and network security' },
  },
  dataMigration: {
    key: 'dataMigration', cost: 2500, icon: 'database',
    label: { es: 'Migración de datos', en: 'Data migration' },
    description: { es: 'Transfiere datos de un sistema a otro', en: 'Transfer data from one system to another' },
  },
  staffTraining: {
    key: 'staffTraining', cost: 2000, icon: 'users',
    label: { es: 'Capacitación del personal', en: 'Staff training' },
    description: { es: 'Entrena a tu equipo en las nuevas herramientas', en: 'Train your team on new tools' },
  },
  postSupport: {
    key: 'postSupport', cost: 2500, icon: 'headphones',
    label: { es: 'Soporte técnico post-implementación', en: 'Post-implementation support' },
    description: { es: 'Asistencia técnica después de la entrega', en: 'Technical assistance after delivery' },
  },
  cloudMigration: {
    key: 'cloudMigration', cost: 4000, icon: 'cloud',
    label: { es: 'Migración a la nube', en: 'Cloud migration' },
    description: { es: 'Mueve tus sistemas a AWS, GCP o Azure', en: 'Move your systems to AWS, GCP or Azure' },
  },
  securityAudit: {
    key: 'securityAudit', cost: 3000, icon: 'shield',
    label: { es: 'Auditoría de seguridad', en: 'Security audit' },
    description: { es: 'Identifica vulnerabilidades y riesgos', en: 'Identify vulnerabilities and risks' },
  },

  // ── AI Integration ───────────────────────────────────────
  whatsappBot: {
    key: 'whatsappBot', cost: 0, icon: 'messageCircle',
    label: { es: 'Bot de WhatsApp', en: 'WhatsApp bot' },
    description: { es: 'Atiende clientes 24/7 por WhatsApp', en: 'Serve customers 24/7 via WhatsApp' },
  },
  webChatbot: {
    key: 'webChatbot', cost: 4000, icon: 'bot',
    label: { es: 'Chatbot para sitio web', en: 'Website chatbot' },
    description: { es: 'Asistente inteligente en tu página web', en: 'Smart assistant on your website' },
  },
  processAutomation: {
    key: 'processAutomation', cost: 5000, icon: 'settings',
    label: { es: 'Automatización de procesos', en: 'Process automation' },
    description: { es: 'Automatiza tareas repetitivas con IA', en: 'Automate repetitive tasks with AI' },
  },
  dataAnalysis: {
    key: 'dataAnalysis', cost: 6000, icon: 'barChart',
    label: { es: 'Análisis de datos con IA', en: 'AI data analysis' },
    description: { es: 'Extrae insights y predicciones de tus datos', en: 'Extract insights and predictions from your data' },
  },
  smartReports: {
    key: 'smartReports', cost: 4000, icon: 'barChart',
    label: { es: 'Reportes inteligentes', en: 'Smart reports' },
    description: { es: 'Reportes generados automáticamente con IA', en: 'AI-generated automatic reports' },
  },
  virtualAssistant: {
    key: 'virtualAssistant', cost: 7500, icon: 'bot',
    label: { es: 'Asistente virtual personalizado', en: 'Custom virtual assistant' },
    description: { es: 'IA entrenada con datos de tu negocio', en: 'AI trained with your business data' },
  },
  voiceAssistant: {
    key: 'voiceAssistant', cost: 6000, icon: 'mic',
    label: { es: 'Asistente de voz', en: 'Voice assistant' },
    description: { es: 'Interacción por voz con IA conversacional', en: 'Voice interaction with conversational AI' },
  },
  docProcessing: {
    key: 'docProcessing', cost: 4000, icon: 'fileText',
    label: { es: 'Procesamiento de documentos', en: 'Document processing' },
    description: { es: 'Extrae datos de facturas, contratos, recibos', en: 'Extract data from invoices, contracts, receipts' },
  },

  // ── Consulting & Architecture ──────────────────────────
  techAudit: {
    key: 'techAudit', cost: 0, icon: 'search',
    label: { es: 'Auditoría técnica', en: 'Technical audit' },
    description: { es: 'Evaluación completa de tu stack y código actual', en: 'Complete assessment of your current stack and code' },
  },
  archDesign: {
    key: 'archDesign', cost: 3000, icon: 'layers',
    label: { es: 'Diseño de arquitectura', en: 'Architecture design' },
    description: { es: 'Diagramas y diseño de arquitectura escalable', en: 'Diagrams and scalable architecture design' },
  },
  codeReview: {
    key: 'codeReview', cost: 2500, icon: 'code',
    label: { es: 'Revisión de código', en: 'Code review' },
    description: { es: 'Revisión profunda de calidad, seguridad y rendimiento', en: 'Deep review of quality, security and performance' },
  },
  roadmap: {
    key: 'roadmap', cost: 2500, icon: 'map',
    label: { es: 'Roadmap tecnológico', en: 'Technology roadmap' },
    description: { es: 'Plan de evolución tecnológica a corto y largo plazo', en: 'Short and long-term technology evolution plan' },
  },
  stackSelection: {
    key: 'stackSelection', cost: 2000, icon: 'settings',
    label: { es: 'Selección de stack', en: 'Stack selection' },
    description: { es: 'Recomendación de tecnologías ideales para tu proyecto', en: 'Ideal technology recommendations for your project' },
  },
  perfOptimization: {
    key: 'perfOptimization', cost: 3500, icon: 'zap',
    label: { es: 'Optimización de rendimiento', en: 'Performance optimization' },
    description: { es: 'Análisis y mejora de tiempos de respuesta', en: 'Analysis and improvement of response times' },
  },
  scalabilityPlan: {
    key: 'scalabilityPlan', cost: 3000, icon: 'trending',
    label: { es: 'Plan de escalabilidad', en: 'Scalability plan' },
    description: { es: 'Estrategia para crecer sin reescribir tu sistema', en: 'Strategy to grow without rewriting your system' },
  },
  docAndDiagrams: {
    key: 'docAndDiagrams', cost: 2000, icon: 'fileText',
    label: { es: 'Documentación y diagramas', en: 'Documentation & diagrams' },
    description: { es: 'Documentación técnica, diagramas de flujo y ERDs', en: 'Technical docs, flow diagrams and ERDs' },
  },

  // ── Team Training ──────────────────────────────────────
  needsAssessment: {
    key: 'needsAssessment', cost: 0, icon: 'clipboard',
    label: { es: 'Evaluación de necesidades', en: 'Needs assessment' },
    description: { es: 'Diagnóstico del nivel y necesidades de tu equipo', en: 'Assessment of your team\'s level and needs' },
  },
  customCurriculum: {
    key: 'customCurriculum', cost: 2000, icon: 'fileText',
    label: { es: 'Currículo personalizado', en: 'Custom curriculum' },
    description: { es: 'Plan de estudios adaptado a tu stack y objetivos', en: 'Study plan adapted to your stack and goals' },
  },
  liveWorkshops: {
    key: 'liveWorkshops', cost: 3000, icon: 'users',
    label: { es: 'Talleres en vivo', en: 'Live workshops' },
    description: { es: 'Sesiones presenciales o remotas con ejercicios prácticos', en: 'On-site or remote sessions with hands-on exercises' },
  },
  trainingMaterials: {
    key: 'trainingMaterials', cost: 1500, icon: 'book',
    label: { es: 'Material de capacitación', en: 'Training materials' },
    description: { es: 'Guías, presentaciones y recursos de referencia', en: 'Guides, presentations and reference resources' },
  },
  practiceProjects: {
    key: 'practiceProjects', cost: 2500, icon: 'code',
    label: { es: 'Proyectos de práctica', en: 'Practice projects' },
    description: { es: 'Ejercicios aplicados al contexto real de tu empresa', en: 'Exercises applied to your company\'s real context' },
  },
  postTrainingSupport: {
    key: 'postTrainingSupport', cost: 2000, icon: 'headphones',
    label: { es: 'Soporte post-capacitación', en: 'Post-training support' },
    description: { es: 'Acompañamiento y resolución de dudas después del curso', en: 'Follow-up and Q&A support after the course' },
  },
  certificationPath: {
    key: 'certificationPath', cost: 1500, icon: 'award',
    label: { es: 'Ruta de certificación', en: 'Certification path' },
    description: { es: 'Evaluaciones y constancias de capacitación', en: 'Assessments and training certificates' },
  },
  recordedSessions: {
    key: 'recordedSessions', cost: 2000, icon: 'video',
    label: { es: 'Sesiones grabadas', en: 'Recorded sessions' },
    description: { es: 'Grabaciones de todas las sesiones para consulta futura', en: 'Recordings of all sessions for future reference' },
  },

  // ── Migration & Modernization ──────────────────────────
  legacyAudit: {
    key: 'legacyAudit', cost: 0, icon: 'search',
    label: { es: 'Auditoría de sistema legacy', en: 'Legacy system audit' },
    description: { es: 'Evaluación completa del sistema actual y sus dependencias', en: 'Complete assessment of current system and dependencies' },
  },
  codeRefactor: {
    key: 'codeRefactor', cost: 5000, icon: 'code',
    label: { es: 'Refactorización de código', en: 'Code refactoring' },
    description: { es: 'Moderniza la estructura y calidad del código existente', en: 'Modernize the structure and quality of existing code' },
  },
  dbMigration: {
    key: 'dbMigration', cost: 4000, icon: 'database',
    label: { es: 'Migración de base de datos', en: 'Database migration' },
    description: { es: 'Migra datos entre motores o versiones de BD', en: 'Migrate data between database engines or versions' },
  },
  cloudMigrationMod: {
    key: 'cloudMigrationMod', cost: 4000, icon: 'cloud',
    label: { es: 'Migración a la nube', en: 'Cloud migration' },
    description: { es: 'Mueve tu infraestructura on-premise a AWS, GCP o Azure', en: 'Move your on-premise infrastructure to AWS, GCP or Azure' },
  },
  apiModernization: {
    key: 'apiModernization', cost: 3500, icon: 'link',
    label: { es: 'Modernización de APIs', en: 'API modernization' },
    description: { es: 'Actualiza APIs monolíticas a microservicios REST/GraphQL', en: 'Update monolithic APIs to REST/GraphQL microservices' },
  },
  testingSetup: {
    key: 'testingSetup', cost: 3000, icon: 'checkCircle',
    label: { es: 'Setup de testing', en: 'Testing setup' },
    description: { es: 'Implementa pruebas unitarias, integración y E2E', en: 'Implement unit, integration and E2E tests' },
  },
  perfTuning: {
    key: 'perfTuning', cost: 3500, icon: 'zap',
    label: { es: 'Optimización de rendimiento', en: 'Performance tuning' },
    description: { es: 'Identifica y elimina cuellos de botella', en: 'Identify and eliminate bottlenecks' },
  },
  documentationMod: {
    key: 'documentationMod', cost: 2000, icon: 'fileText',
    label: { es: 'Documentación del sistema', en: 'System documentation' },
    description: { es: 'Documentación técnica del sistema modernizado', en: 'Technical documentation of the modernized system' },
  },
};

export const BUSINESS_SIZES = [
  { key: '1-5', label: { es: '1-5 empleados', en: '1-5 employees' }, multiplier: 1.0, icon: 'user' },
  { key: '6-20', label: { es: '6-20 empleados', en: '6-20 employees' }, multiplier: 1.15, icon: 'users' },
  { key: '21-50', label: { es: '21-50 empleados', en: '21-50 employees' }, multiplier: 1.3, icon: 'building' },
  { key: '50+', label: { es: '50+ empleados', en: '50+ employees' }, multiplier: 1.5, icon: 'city' },
] as const;

export const CURRENT_STATES = [
  { key: 'fromScratch', label: { es: 'Empezar desde cero', en: 'Start from scratch' }, multiplier: 1.0, icon: 'plus' },
  { key: 'improve', label: { es: 'Mejorar lo que ya tengo', en: 'Improve what I have' }, multiplier: 0.7, icon: 'refresh' },
  { key: 'migrate', label: { es: 'Migrar de otro sistema', en: 'Migrate from another system' }, multiplier: 1.2, icon: 'shuffle' },
] as const;

export const TIMELINES = [
  { key: 'asap', label: { es: 'Lo antes posible', en: 'ASAP' }, multiplier: 1.3, icon: 'zap' },
  { key: '1-3months', label: { es: '1-3 meses', en: '1-3 months' }, multiplier: 1.0, icon: 'calendar' },
  { key: '3-6months', label: { es: '3-6 meses', en: '3-6 months' }, multiplier: 0.95, icon: 'calendarRange' },
  { key: 'exploring', label: { es: 'Solo estoy explorando', en: 'Just exploring' }, multiplier: 1.0, icon: 'search' },
] as const;

export const CURRENCIES = [
  { key: 'MXN', label: 'MXN', flag: '🇲🇽', name: { es: 'Pesos mexicanos', en: 'Mexican Pesos' } },
  { key: 'USD', label: 'USD', flag: '🇺🇸', name: { es: 'Dólares americanos', en: 'US Dollars' } },
] as const;

export const EXCHANGE_RATE = 17.5;

export interface QuoteSelections {
  projectTypes: string[];
  features: string[];
  businessSize: string;
  currentState: string;
  timeline: string;
  currency: string;
}

export interface QuoteResult {
  min: number;
  max: number;
  total: number;
  currency: string;
}

export function calculateQuote(selections: QuoteSelections): QuoteResult {
  // Sum project type bases
  let base = 0;
  for (const key of selections.projectTypes) {
    const pt = PROJECT_TYPES.find((p) => p.key === key);
    if (pt) base += pt.base;
  }

  // Sum feature costs
  let featureCost = 0;
  for (const key of selections.features) {
    const f = FEATURES[key];
    if (f) featureCost += f.cost;
  }

  let total = base + featureCost;

  // Apply multipliers
  const sizeEntry = BUSINESS_SIZES.find((s) => s.key === selections.businessSize);
  if (sizeEntry) total *= sizeEntry.multiplier;

  const stateEntry = CURRENT_STATES.find((s) => s.key === selections.currentState);
  if (stateEntry) total *= stateEntry.multiplier;

  const timelineEntry = TIMELINES.find((t) => t.key === selections.timeline);
  if (timelineEntry) total *= timelineEntry.multiplier;

  // Convert currency if USD
  if (selections.currency === 'USD') {
    total = total / EXCHANGE_RATE;
  }

  return {
    min: Math.round(total * 0.85),
    max: Math.round(total * 1.15),
    total: Math.round(total),
    currency: selections.currency,
  };
}

function formatPlanPrice(n: number, currency: string, lang: 'es' | 'en'): string {
  return new Intl.NumberFormat(lang === 'es' ? 'es-MX' : 'en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(Math.round(n));
}

const HOURLY_RATE: Record<string, number> = { MXN: 500, USD: 30 };
const SERVER_COST: Record<string, number> = { MXN: 3000, USD: 170 };

export const PAYMENT_PLANS: PaymentPlan[] = [
  {
    key: 'fullPayment',
    label: { es: 'Pago completo', en: 'Full payment' },
    description: { es: 'Un solo pago con descuento', en: 'Single payment with discount' },
    icon: 'dollarSign',
    badge: { es: 'Ahorro 10%', en: 'Save 10%' },
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const discounted = total * 0.90;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon, badge: this.badge,
        primary: formatPlanPrice(discounted, currency, lang),
        secondary: lang === 'es' ? 'Pago único' : 'One-time payment',
        details: {
          es: 'Paga el total al inicio y obtén un 10% de descuento.',
          en: 'Pay the full amount upfront and get a 10% discount.',
        },
        totalCost: discounted,
      };
    },
  },
  {
    key: 'splitPayment',
    label: { es: '50% / 50%', en: '50% / 50%' },
    description: { es: 'Mitad al inicio, mitad al entregar', en: 'Half upfront, half on delivery' },
    icon: 'creditCard',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const half = total / 2;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: formatPlanPrice(half, currency, lang),
        secondary: lang === 'es' ? '2 pagos' : '2 payments',
        details: {
          es: '50% al iniciar el proyecto y 50% al entregarlo.',
          en: '50% when the project starts and 50% on delivery.',
        },
        totalCost: total,
      };
    },
  },
  {
    key: 'msi3',
    label: { es: '3 MSI', en: '3 installments' },
    description: { es: '3 meses sin intereses', en: '3 months interest-free' },
    icon: 'calendar',
    badge: { es: 'Popular', en: 'Popular' },
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const monthly = total / 3;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon, badge: this.badge,
        primary: formatPlanPrice(monthly, currency, lang),
        secondary: lang === 'es' ? '× 3 meses' : '× 3 months',
        details: {
          es: 'Divide el costo en 3 pagos mensuales sin intereses.',
          en: 'Split the cost into 3 monthly interest-free payments.',
        },
        totalCost: total,
      };
    },
  },
  {
    key: 'msi6',
    label: { es: '6 MSI', en: '6 installments' },
    description: { es: '6 meses sin intereses', en: '6 months interest-free' },
    icon: 'calendar',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const monthly = total / 6;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: formatPlanPrice(monthly, currency, lang),
        secondary: lang === 'es' ? '× 6 meses' : '× 6 months',
        details: {
          es: 'Divide el costo en 6 pagos mensuales sin intereses.',
          en: 'Split the cost into 6 monthly interest-free payments.',
        },
        totalCost: total,
      };
    },
  },
  {
    key: 'financing12',
    label: { es: '12 meses', en: '12 months' },
    description: { es: 'Financiamiento a 12 meses', en: '12-month financing' },
    icon: 'clock',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const financed = total * 1.15;
      const monthly = financed / 12;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: formatPlanPrice(monthly, currency, lang),
        secondary: lang === 'es' ? '× 12 meses (+15%)' : '× 12 months (+15%)',
        details: {
          es: 'Pagos mensuales a 12 meses con un 15% de recargo por financiamiento.',
          en: 'Monthly payments over 12 months with a 15% financing surcharge.',
        },
        totalCost: financed,
      };
    },
  },
  {
    key: 'saasMonthly',
    label: { es: 'SaaS mensual', en: 'Monthly SaaS' },
    description: { es: 'Pago mensual con mantenimiento', en: 'Monthly payment with maintenance' },
    icon: 'cloud',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const monthly = (total + total * 0.15) / 12;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: formatPlanPrice(monthly, currency, lang) + (lang === 'es' ? '/mes' : '/mo'),
        secondary: lang === 'es' ? '× 12 meses (incluye mantenimiento)' : '× 12 months (includes maintenance)',
        details: {
          es: 'Pago mensual que incluye desarrollo + mantenimiento y actualizaciones.',
          en: 'Monthly payment that includes development + maintenance and updates.',
        },
        totalCost: total + total * 0.15,
      };
    },
  },
  {
    key: 'annualLicense',
    label: { es: 'Licencia anual', en: 'Annual license' },
    description: { es: 'Licencia + servidor mensual', en: 'License + monthly server' },
    icon: 'server',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const upfront = total * 0.60;
      const server = SERVER_COST[currency] ?? SERVER_COST.MXN;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: formatPlanPrice(upfront, currency, lang),
        secondary: `+ ${formatPlanPrice(server, currency, lang)}/${lang === 'es' ? 'mes servidor' : 'mo server'}`,
        details: {
          es: `60% del costo como licencia inicial + ${formatPlanPrice(server, currency, lang)}/mes por servidor y soporte.`,
          en: `60% of the cost as initial license + ${formatPlanPrice(server, currency, lang)}/mo for server and support.`,
        },
        totalCost: upfront + server * 12,
      };
    },
  },
  {
    key: 'timeRetainer',
    label: { es: 'Por horas', en: 'Hourly retainer' },
    description: { es: 'Retainer de horas estimadas', en: 'Estimated hours retainer' },
    icon: 'clock',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const rate = HOURLY_RATE[currency] ?? HOURLY_RATE.MXN;
      const hours = Math.ceil(total / rate);
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: `~${hours} hrs`,
        secondary: `@ ${formatPlanPrice(rate, currency, lang)}/hr`,
        details: {
          es: `Aproximadamente ${hours} horas de desarrollo a ${formatPlanPrice(rate, currency, lang)} por hora.`,
          en: `Approximately ${hours} development hours at ${formatPlanPrice(rate, currency, lang)} per hour.`,
        },
        totalCost: hours * rate,
      };
    },
  },
  {
    key: 'payroll',
    label: { es: 'Nómina', en: 'Payroll-style' },
    description: { es: 'Equivalente a nómina mensual', en: 'Monthly payroll equivalent' },
    icon: 'users',
    calculate(total, currency) {
      const lang: 'es' | 'en' = currency === 'USD' ? 'en' : 'es';
      const monthly = total / 6;
      return {
        key: this.key, label: this.label, description: this.description,
        icon: this.icon,
        primary: formatPlanPrice(monthly, currency, lang) + (lang === 'es' ? '/mes' : '/mo'),
        secondary: lang === 'es' ? '× 6 meses' : '× 6 months',
        details: {
          es: 'Paga como si fuera un salario mensual durante 6 meses de desarrollo dedicado.',
          en: 'Pay as a monthly salary over 6 months of dedicated development.',
        },
        totalCost: total,
      };
    },
  },
];

export function generatePaymentPlans(
  total: number,
  currency: string,
  _lang: 'es' | 'en',
): GeneratedPlan[] {
  return PAYMENT_PLANS.map((plan) => plan.calculate(total, currency));
}
