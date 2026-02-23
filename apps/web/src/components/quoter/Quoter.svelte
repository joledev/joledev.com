<script lang="ts">
  import { toast } from '../../lib/toast.svelte';
  import {
    PROJECT_TYPES,
    FEATURES,
    BUSINESS_SIZES,
    CURRENT_STATES,
    TIMELINES,
    CURRENCIES,
    calculateQuote,
    generatePaymentPlans,
    SOURCE_CODE_SURCHARGE,
    type QuoteSelections,
    type QuoteResult,
    type GeneratedPlan,
  } from '../../lib/quoter-config';

  interface Props {
    lang: 'es' | 'en';
    apiUrl?: string;
    turnstileSiteKey?: string;
  }

  let { lang, apiUrl = '', turnstileSiteKey = '' }: Props = $props();

  const TOTAL_STEPS = 7;
  let currentStep = $state(1);

  // Selections state
  let selectedProjectTypes = $state<string[]>([]);
  let selectedFeatures = $state<string[]>([]);
  let businessSize = $state('');
  let currentState = $state('');
  let timeline = $state('');
  let currency = $state('');

  // Panel ref for scroll-to-top on step change
  let panelEl = $state<HTMLElement | null>(null);

  function scrollToPanel() {
    panelEl?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  // Step 1 hover/click detail
  let hoveredProjectKey = $state<string | null>(null);
  let tappedDetailKey = $state<string | null>(null);

  function handleProjectEnter(key: string) {
    tappedDetailKey = null;
    hoveredProjectKey = key;
  }

  function handleProjectLeave() {
    hoveredProjectKey = null;
  }

  // Step 2 collapsible groups
  let collapsedGroups = $state<Set<string>>(new Set());

  function toggleGroup(key: string) {
    const next = new Set(collapsedGroups);
    if (next.has(key)) {
      next.delete(key);
    } else {
      next.add(key);
    }
    collapsedGroups = next;
  }

  // Contact form
  let contactName = $state('');
  let contactEmail = $state('');
  let contactPhone = $state('');
  let contactCompany = $state('');
  let contactNotes = $state('');
  let formErrors = $state<Record<string, string>>({});
  let fieldTouched = $state<Record<string, boolean>>({});

  const EMAIL_RE = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$/;

  function validateField(field: string) {
    const errors = { ...formErrors };
    if (field === 'name') {
      if (!contactName.trim()) {
        errors.name = lang === 'es' ? 'El nombre es requerido' : 'Name is required';
      } else if (contactName.trim().length > 200) {
        errors.name = lang === 'es' ? 'Máximo 200 caracteres' : 'Maximum 200 characters';
      } else {
        delete errors.name;
      }
    }
    if (field === 'email') {
      if (!contactEmail.trim()) {
        errors.email = lang === 'es' ? 'El email es requerido' : 'Email is required';
      } else if (!EMAIL_RE.test(contactEmail)) {
        errors.email = lang === 'es' ? 'Email no válido' : 'Invalid email';
      } else {
        delete errors.email;
      }
    }
    formErrors = errors;
  }

  // Payment plan state
  let selectedPlanKey = $state('fullPayment');
  let includeSourceCode = $state(false);
  let showIncludes = $state(false);

  // Submission
  let submitting = $state(false);
  let submitted = $state(false);
  let submitError = $state('');

  // Turnstile CAPTCHA
  let turnstileToken = $state('');
  let turnstileWidgetId = $state<string | null>(null);

  function initTurnstile(container: HTMLElement) {
    if (!turnstileSiteKey || turnstileWidgetId) return;
    if (!(window as any).turnstile) return;
    turnstileWidgetId = (window as any).turnstile.render(container, {
      sitekey: turnstileSiteKey,
      callback: (token: string) => { turnstileToken = token; },
      'expired-callback': () => { turnstileToken = ''; },
      'error-callback': () => { turnstileToken = ''; },
      theme: 'auto',
      size: 'flexible',
    });
  }

  function resetTurnstile() {
    if (turnstileWidgetId && (window as any).turnstile) {
      (window as any).turnstile.reset(turnstileWidgetId);
      turnstileToken = '';
    }
  }

  // Load Turnstile script dynamically
  $effect(() => {
    if (!turnstileSiteKey) return;
    if (document.querySelector('script[src*="challenges.cloudflare.com/turnstile"]')) return;
    const script = document.createElement('script');
    script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit';
    script.async = true;
    script.defer = true;
    document.head.appendChild(script);
  });

  // Derived
  let availableFeatures = $derived(() => {
    const featureKeys = new Set<string>();
    for (const key of selectedProjectTypes) {
      const pt = PROJECT_TYPES.find((p) => p.key === key);
      if (pt) pt.features.forEach((f) => featureKeys.add(f));
    }
    return Array.from(featureKeys);
  });

  let quoteResult = $derived(() => {
    if (!businessSize || !currentState || !timeline || !currency) return null;
    return calculateQuote({
      projectTypes: selectedProjectTypes,
      features: selectedFeatures,
      businessSize,
      currentState,
      timeline,
      currency,
    });
  });

  let adjustedTotal = $derived(() => {
    const result = quoteResult();
    if (!result) return 0;
    return includeSourceCode ? result.total * (1 + SOURCE_CODE_SURCHARGE) : result.total;
  });

  let paymentPlans = $derived(() => {
    const t = adjustedTotal();
    const result = quoteResult();
    if (!t || !result) return [];
    return generatePaymentPlans(t, result.currency, lang);
  });

  let canNext = $derived(() => {
    switch (currentStep) {
      case 1: return selectedProjectTypes.length > 0;
      case 2: return true; // features are optional
      case 3: return businessSize !== '';
      case 4: return currentState !== '';
      case 5: return timeline !== '';
      case 6: return currency !== '';
      case 7: return true;
      default: return false;
    }
  });

  let selectedSummary = $derived(() => {
    const items: string[] = [];
    for (const key of selectedProjectTypes) {
      const pt = PROJECT_TYPES.find((p) => p.key === key);
      if (pt) items.push(pt.label[lang]);
    }
    for (const key of selectedFeatures) {
      const f = FEATURES[key];
      if (f) items.push(f.label[lang]);
    }
    return items;
  });

  function toggleProjectType(key: string) {
    if (selectedProjectTypes.includes(key)) {
      selectedProjectTypes = selectedProjectTypes.filter((k) => k !== key);
      // Remove features that are no longer available
      const stillAvailable = new Set<string>();
      for (const ptKey of selectedProjectTypes) {
        const pt = PROJECT_TYPES.find((p) => p.key === ptKey);
        if (pt) pt.features.forEach((f) => stillAvailable.add(f));
      }
      selectedFeatures = selectedFeatures.filter((f) => stillAvailable.has(f));
    } else {
      selectedProjectTypes = [...selectedProjectTypes, key];
    }
  }

  function toggleFeature(key: string) {
    if (selectedFeatures.includes(key)) {
      selectedFeatures = selectedFeatures.filter((k) => k !== key);
    } else {
      selectedFeatures = [...selectedFeatures, key];
    }
  }

  function selectedCountForGroup(pt: { features: string[] }): number {
    return pt.features.filter((f) => selectedFeatures.includes(f)).length;
  }

  function nextStep() {
    if (canNext() && currentStep < TOTAL_STEPS) {
      currentStep++;
      scrollToPanel();
    }
  }

  function prevStep() {
    if (currentStep > 1) {
      currentStep--;
      scrollToPanel();
    }
  }

  function formatPrice(n: number, cur: string): string {
    return new Intl.NumberFormat(lang === 'es' ? 'es-MX' : 'en-US', {
      style: 'currency',
      currency: cur,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(n);
  }

  function formatFeaturePrice(cost: number): string {
    if (cost >= 1000) {
      const k = cost / 1000;
      return `+$${k % 1 === 0 ? k.toFixed(0) : k.toFixed(1)}k`;
    }
    return `+$${cost}`;
  }

  function validateForm(): boolean {
    const errors: Record<string, string> = {};
    if (!contactName.trim()) {
      errors.name = lang === 'es' ? 'El nombre es requerido' : 'Name is required';
    } else if (contactName.trim().length > 200) {
      errors.name = lang === 'es' ? 'Máximo 200 caracteres' : 'Maximum 200 characters';
    }
    if (!contactEmail.trim()) {
      errors.email = lang === 'es' ? 'El email es requerido' : 'Email is required';
    } else if (!EMAIL_RE.test(contactEmail)) {
      errors.email = lang === 'es' ? 'Email no válido' : 'Invalid email';
    }
    formErrors = errors;
    return Object.keys(errors).length === 0;
  }

  async function submitQuote() {
    if (!validateForm()) return;

    const result = quoteResult();
    if (!result) return;

    submitting = true;
    submitError = '';

    const body = {
      projectTypes: selectedProjectTypes,
      features: selectedFeatures,
      businessSize,
      currentState,
      timeline,
      currency,
      estimatedMin: result.min,
      estimatedMax: result.max,
      paymentPlan: selectedPlanKey,
      includeSourceCode,
      contact: {
        name: contactName.trim(),
        email: contactEmail.trim(),
        phone: contactPhone.trim(),
        company: contactCompany.trim(),
        notes: contactNotes.trim(),
      },
      lang,
      turnstileToken,
    };

    try {
      const res = await fetch(`${apiUrl}/quotes`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });

      if (!res.ok) {
        const data = await res.json().catch(() => null);
        throw new Error(data?.message || `Error ${res.status}`);
      }

      submitted = true;
      toast.success(labels.successTitle);
    } catch (err: unknown) {
      submitError =
        lang === 'es'
          ? 'Hubo un error al enviar. Por favor intenta de nuevo.'
          : 'There was an error submitting. Please try again.';
      toast.error(submitError);
      resetTurnstile();
    } finally {
      submitting = false;
    }
  }

  // Labels
  const labels = {
    next: lang === 'es' ? 'Siguiente' : 'Next',
    prev: lang === 'es' ? 'Anterior' : 'Previous',
    step: lang === 'es' ? 'Paso' : 'Step',
    of: lang === 'es' ? 'de' : 'of',
    send: lang === 'es' ? 'Enviar cotización' : 'Send Quote',
    sending: lang === 'es' ? 'Enviando...' : 'Sending...',
    schedule: lang === 'es' ? 'Agendar una reunión' : 'Schedule a meeting',
    successTitle: lang === 'es' ? '¡Cotización enviada!' : 'Quote sent!',
    successMsg: lang === 'es' ? 'Te contactaré en las próximas 24 horas.' : "I'll contact you within 24 hours.",
    newQuote: lang === 'es' ? 'Nueva cotización' : 'New quote',
    stepTitles: [
      '',
      lang === 'es' ? '¿Qué tipo de proyecto necesitas?' : 'What type of project do you need?',
      lang === 'es' ? '¿Qué funcionalidades necesitas?' : 'What features do you need?',
      lang === 'es' ? '¿Cuál es el tamaño de tu negocio?' : "What's your business size?",
      lang === 'es' ? '¿Cuál es tu situación actual?' : "What's your current situation?",
      lang === 'es' ? '¿Para cuándo lo necesitas?' : 'When do you need it?',
      lang === 'es' ? '¿En qué moneda prefieres tu presupuesto?' : 'Preferred currency?',
      lang === 'es' ? 'Tu presupuesto estimado' : 'Your estimated budget',
    ],
    multiSelect: lang === 'es' ? 'Selección múltiple permitida' : 'Multiple selection allowed',
    featuresOptional: lang === 'es' ? 'Selecciona las que apliquen a tu proyecto' : 'Select the ones that apply to your project',
    estimate: lang === 'es' ? 'Tu presupuesto estimado' : 'Your estimated budget',
    includes: lang === 'es' ? 'Esto incluye:' : 'This includes:',
    contactTitle: lang === 'es' ? 'Para recibir una propuesta detallada:' : 'To receive a detailed proposal:',
    name: lang === 'es' ? 'Nombre' : 'Name',
    email: 'Email',
    phone: lang === 'es' ? 'Teléfono' : 'Phone',
    company: lang === 'es' ? 'Empresa' : 'Company',
    notes: lang === 'es' ? 'Notas adicionales' : 'Additional notes',
    included: lang === 'es' ? 'Incluido' : 'Included',
    sourceCodeLabel: lang === 'es' ? 'Incluir código fuente' : 'Include source code',
    sourceCodeNote: lang === 'es' ? '+25% — Recibes todo el código del proyecto' : '+25% — You receive all project source code',
    choosePlan: lang === 'es' ? 'Elige tu plan de pago' : 'Choose your payment plan',
    baseBudget: lang === 'es' ? 'Presupuesto base' : 'Base budget',
  };

  // SVG icons map
  const icons: Record<string, string> = {
    // Project type icons (original)
    globe: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/></svg>',
    monitor: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>',
    package: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg>',
    dollarSign: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>',
    fileText: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>',
    refreshCw: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>',
    bot: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="10" rx="2"/><circle cx="12" cy="5" r="2"/><path d="M12 7v4"/><line x1="8" y1="16" x2="8" y2="16"/><line x1="16" y1="16" x2="16" y2="16"/></svg>',

    // New project type icons
    shoppingCart: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>',
    smartphone: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"/><line x1="12" y1="18" x2="12.01" y2="18"/></svg>',
    cloud: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"/></svg>',
    calendarCheck: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><path d="M9 16l2 2 4-4"/></svg>',
    link: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>',
    server: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"/><rect x="2" y="14" width="20" height="8" rx="2" ry="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>',

    // Feature icons
    mail: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>',
    search: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>',
    settings: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>',
    share: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/></svg>',
    sparkles: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3l1.912 5.813a2 2 0 0 0 1.275 1.275L21 12l-5.813 1.912a2 2 0 0 0-1.275 1.275L12 21l-1.912-5.813a2 2 0 0 0-1.275-1.275L3 12l5.813-1.912a2 2 0 0 0 1.275-1.275L12 3z"/></svg>',
    barChart: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="20" x2="12" y2="10"/><line x1="18" y1="20" x2="18" y2="4"/><line x1="6" y1="20" x2="6" y2="16"/></svg>',
    messageCircle: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z"/></svg>',
    grid: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>',
    creditCard: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="1" y="4" width="22" height="16" rx="2" ry="2"/><line x1="1" y1="10" x2="23" y2="10"/></svg>',
    truck: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="1" y="3" width="15" height="13"/><polygon points="16 8 20 8 23 11 23 16 16 16 16 8"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>',
    tag: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/><line x1="7" y1="7" x2="7.01" y2="7"/></svg>',
    star: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>',
    heart: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/></svg>',
    mapPin: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>',
    bell: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>',
    wifi: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12.55a11 11 0 0 1 14.08 0"/><path d="M1.42 9a16 16 0 0 1 21.16 0"/><path d="M8.53 16.11a6 6 0 0 1 6.95 0"/><line x1="12" y1="20" x2="12.01" y2="20"/></svg>',
    shield: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>',
    camera: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/><circle cx="12" cy="13" r="4"/></svg>',
    upload: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/><path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/></svg>',
    users: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>',
    download: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>',
    building: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="2" width="16" height="20" rx="2" ry="2"/><line x1="9" y1="6" x2="9.01" y2="6"/><line x1="15" y1="6" x2="15.01" y2="6"/><line x1="9" y1="10" x2="9.01" y2="10"/><line x1="15" y1="10" x2="15.01" y2="10"/><line x1="9" y1="14" x2="9.01" y2="14"/><line x1="15" y1="14" x2="15.01" y2="14"/><line x1="9" y1="18" x2="15" y2="18"/></svg>',
    clipboard: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1" ry="1"/></svg>',
    gitBranch: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="6" y1="3" x2="6" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/></svg>',
    layers: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>',
    userPlus: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="8.5" cy="7" r="4"/><line x1="20" y1="8" x2="20" y2="14"/><line x1="23" y1="11" x2="17" y2="11"/></svg>',
    code: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>',
    lock: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>',
    palette: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="13.5" cy="6.5" r="2.5"/><circle cx="19" cy="11.5" r="2.5"/><circle cx="6" cy="12.5" r="2.5"/><circle cx="10" cy="18.5" r="2.5"/><path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10c.93 0 1.5-.67 1.5-1.5 0-.39-.15-.74-.39-1.04-.23-.29-.38-.63-.38-1.04 0-.83.67-1.5 1.5-1.5H16c3.31 0 6-2.69 6-6 0-5.17-4.49-9-10-9z"/></svg>',
    scan: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7V5a2 2 0 0 1 2-2h2"/><path d="M17 3h2a2 2 0 0 1 2 2v2"/><path d="M21 17v2a2 2 0 0 1-2 2h-2"/><path d="M7 21H5a2 2 0 0 1-2-2v-2"/><line x1="7" y1="12" x2="17" y2="12"/></svg>',
    award: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="8" r="7"/><polyline points="8.21 13.89 7 23 12 20 17 23 15.79 13.88"/></svg>',
    database: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>',
    calculator: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="2" width="16" height="20" rx="2"/><line x1="8" y1="6" x2="16" y2="6"/><line x1="16" y1="14" x2="16" y2="18"/><line x1="8" y1="11" x2="8.01" y2="11"/><line x1="12" y1="11" x2="12.01" y2="11"/><line x1="16" y1="11" x2="16.01" y2="11"/><line x1="8" y1="15" x2="8.01" y2="15"/><line x1="12" y1="15" x2="12.01" y2="15"/><line x1="8" y1="19" x2="8.01" y2="19"/><line x1="12" y1="19" x2="12.01" y2="19"/></svg>',
    files: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M10 12H6"/><path d="M14 12h-2"/></svg>',
    calendar: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>',
    clock: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>',
    box: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/></svg>',
    terminal: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>',
    activity: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/></svg>',
    trending: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>',
    zap: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/></svg>',
    hardDrive: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="12" x2="2" y2="12"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/><line x1="6" y1="16" x2="6.01" y2="16"/><line x1="10" y1="16" x2="10.01" y2="16"/></svg>',
    headphones: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 18v-6a9 9 0 0 1 18 0v6"/><path d="M21 19a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3zM3 19a2 2 0 0 0 2 2h1a2 2 0 0 0 2-2v-3a2 2 0 0 0-2-2H3z"/></svg>',
    mic: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"/><path d="M19 10v2a7 7 0 0 1-14 0v-2"/><line x1="12" y1="19" x2="12" y2="23"/><line x1="8" y1="23" x2="16" y2="23"/></svg>',
    chevronDown: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>',
    compass: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polygon points="16.24 7.76 14.12 14.12 7.76 16.24 9.88 9.88 16.24 7.76"/></svg>',
    graduationCap: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 10v6M2 10l10-5 10 5-10 5z"/><path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5"/></svg>',
    arrowUpRight: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="7" y1="17" x2="17" y2="7"/><polyline points="7 7 17 7 17 17"/></svg>',
    book: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>',
    video: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>',
    map: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="1 6 1 22 8 18 16 22 23 18 23 2 16 6 8 2 1 6"/><line x1="8" y1="2" x2="8" y2="18"/><line x1="16" y1="6" x2="16" y2="22"/></svg>',
    checkCircle: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
    calendarRange: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/><line x1="8" y1="14" x2="16" y2="14"/></svg>',
    city: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="1" y="6" width="7" height="16"/><rect x="10" y="2" width="7" height="20"/><rect x="19" y="9" width="4" height="13"/><line x1="4" y1="10" x2="4" y2="10.01"/><line x1="4" y1="14" x2="4" y2="14.01"/><line x1="13" y1="6" x2="13" y2="6.01"/><line x1="13" y1="10" x2="13" y2="10.01"/><line x1="13" y1="14" x2="13" y2="14.01"/><line x1="21" y1="13" x2="21" y2="13.01"/></svg>',
    user: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>',
    plus: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>',
    refresh: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>',
    shuffle: '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 3 21 3 21 8"/><line x1="4" y1="20" x2="21" y2="3"/><polyline points="21 16 21 21 16 21"/><line x1="15" y1="15" x2="21" y2="21"/><line x1="4" y1="4" x2="9" y2="9"/></svg>',
  };
</script>

{#if submitted}
  <!-- Success screen -->
  <div class="quoter-panel">
    <div class="success-screen">
      <div class="success-pulse"></div>
      <div class="success-check">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
      </div>
      <h2 class="success-title">{labels.successTitle}</h2>
      <p class="success-msg">{labels.successMsg}</p>
      <div class="success-actions">
        <a class="quoter-nav-btn next" href={lang === 'es' ? '/es/agendar/' : '/en/schedule/'}>
          {labels.schedule}
        </a>
        <button class="quoter-nav-btn back" onclick={() => { submitted = false; currentStep = 1; selectedProjectTypes = []; selectedFeatures = []; businessSize = ''; currentState = ''; timeline = ''; currency = ''; contactName = ''; contactEmail = ''; contactPhone = ''; contactCompany = ''; contactNotes = ''; }} type="button">
          {labels.newQuote}
        </button>
      </div>
    </div>
  </div>
{:else}
  <div class="quoter-panel" bind:this={panelEl}>
    <!-- Progress bar -->
    <div class="progress-bar">
      <div class="progress-fill" style="width: {(currentStep / TOTAL_STEPS) * 100}%"></div>
    </div>
    <p class="step-indicator">{labels.step} {currentStep} {labels.of} {TOTAL_STEPS}</p>

    <!-- Step title -->
    <h2 class="step-title">{labels.stepTitles[currentStep]}</h2>
    {#if currentStep === 1}
      <p class="step-hint">{labels.multiSelect}</p>
    {/if}
    {#if currentStep === 2}
      <p class="step-hint">{labels.featuresOptional}</p>
    {/if}

    <!-- Step content -->
    {#key currentStep}
    <div class="step-content">
      {#if currentStep === 1}
        <div class="option-grid">
          {#each PROJECT_TYPES as pt, i}
            <button
              class="option-card"
              class:selected={selectedProjectTypes.includes(pt.key)}
              class:expanded={hoveredProjectKey === pt.key || tappedDetailKey === pt.key}
              onclick={() => { toggleProjectType(pt.key); tappedDetailKey = pt.key; }}
              onmouseenter={() => handleProjectEnter(pt.key)}
              onmouseleave={handleProjectLeave}
              onfocus={() => handleProjectEnter(pt.key)}
              type="button"
              aria-pressed={selectedProjectTypes.includes(pt.key)}
              style="animation-delay: {i * 30}ms"
            >
              <span class="option-icon">{@html icons[pt.icon] ?? ''}</span>
              <span class="option-label">{pt.label[lang]}</span>
              {#if selectedProjectTypes.includes(pt.key)}
                <span class="check-mark">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                </span>
              {/if}
              {#if hoveredProjectKey === pt.key || tappedDetailKey === pt.key}
                <span class="option-detail">
                  {#if pt.image}
                    <img
                      class="option-detail-img"
                      src={pt.image}
                      alt={pt.label[lang]}
                      width="280"
                      height="140"
                      loading="lazy"
                    />
                  {/if}
                  {#if pt.description}
                    <span class="option-detail-desc">{pt.description[lang]}</span>
                  {/if}
                </span>
              {/if}
            </button>
          {/each}
        </div>

      {:else if currentStep === 2}
        <div class="features-list">
          {#each selectedProjectTypes as ptKey}
            {@const pt = PROJECT_TYPES.find(p => p.key === ptKey)}
            {#if pt}
              <div class="feature-group">
                <button
                  class="feature-group-header"
                  onclick={() => toggleGroup(pt.key)}
                  type="button"
                  aria-expanded={!collapsedGroups.has(pt.key)}
                >
                  <span class="feature-group-icon">{@html icons[pt.icon] ?? ''}</span>
                  <span class="feature-group-title">{pt.label[lang]}</span>
                  <span class="feature-group-count">{selectedCountForGroup(pt)}/{pt.features.length}</span>
                  <span class="feature-group-chevron" class:collapsed={collapsedGroups.has(pt.key)}>
                    {@html icons.chevronDown}
                  </span>
                </button>
                {#if !collapsedGroups.has(pt.key)}
                  <div class="feature-cards-grid">
                    {#each pt.features as fKey}
                      {@const feature = FEATURES[fKey]}
                      {#if feature}
                        <button
                          class="feature-card"
                          class:selected={selectedFeatures.includes(fKey)}
                          onclick={() => toggleFeature(fKey)}
                          type="button"
                          aria-pressed={selectedFeatures.includes(fKey)}
                        >
                          <div class="feature-card-top">
                            <span class="feature-card-icon">
                              {@html icons[feature.icon ?? ''] ?? ''}
                            </span>
                            {#if feature.cost === 0}
                              <span class="feature-badge included">{labels.included}</span>
                            {:else}
                              <span class="feature-badge price">{formatFeaturePrice(feature.cost)}</span>
                            {/if}
                          </div>
                          <span class="feature-card-name">{feature.label[lang]}</span>
                          {#if feature.description}
                            <span class="feature-card-desc">{feature.description[lang]}</span>
                          {/if}
                          {#if selectedFeatures.includes(fKey)}
                            <span class="check-mark">
                              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                            </span>
                          {/if}
                        </button>
                      {/if}
                    {/each}
                  </div>
                {/if}
              </div>
            {/if}
          {/each}
        </div>

      {:else if currentStep === 3}
        <div class="option-grid cols-2">
          {#each BUSINESS_SIZES as size, i}
            <button
              class="option-card"
              class:selected={businessSize === size.key}
              onclick={() => (businessSize = size.key)}
              type="button"
              aria-pressed={businessSize === size.key}
              style="animation-delay: {i * 30}ms"
            >
              <span class="option-label">{size.label[lang]}</span>
            </button>
          {/each}
        </div>

      {:else if currentStep === 4}
        <div class="option-grid cols-3">
          {#each CURRENT_STATES as state, i}
            <button
              class="option-card"
              class:selected={currentState === state.key}
              onclick={() => (currentState = state.key)}
              type="button"
              aria-pressed={currentState === state.key}
              style="animation-delay: {i * 30}ms"
            >
              <span class="option-label">{state.label[lang]}</span>
            </button>
          {/each}
        </div>

      {:else if currentStep === 5}
        <div class="option-grid cols-2">
          {#each TIMELINES as tl, i}
            <button
              class="option-card"
              class:selected={timeline === tl.key}
              onclick={() => (timeline = tl.key)}
              type="button"
              aria-pressed={timeline === tl.key}
              style="animation-delay: {i * 30}ms"
            >
              <span class="option-label">{tl.label[lang]}</span>
            </button>
          {/each}
        </div>

      {:else if currentStep === 6}
        <div class="option-grid cols-2">
          {#each CURRENCIES as cur, i}
            <button
              class="option-card currency-card"
              class:selected={currency === cur.key}
              onclick={() => (currency = cur.key)}
              type="button"
              aria-pressed={currency === cur.key}
              style="animation-delay: {i * 30}ms"
            >
              <span class="currency-flag">{cur.flag}</span>
              <span class="option-label">{cur.label}</span>
              <span class="currency-name">{cur.name[lang]}</span>
            </button>
          {/each}
        </div>

      {:else if currentStep === 7}
        {@const result = quoteResult()}
        {@const plans = paymentPlans()}
        {@const surchargeMultiplier = includeSourceCode ? 1 + SOURCE_CODE_SURCHARGE : 1}
        <div class="result-section">
          {#if result}
            <div class="estimate-display">
              <p class="estimate-label">{labels.baseBudget}</p>
              <p class="estimate-range">
                {formatPrice(result.min * surchargeMultiplier, result.currency)} — {formatPrice(result.max * surchargeMultiplier, result.currency)}
              </p>
            </div>

            <label class="source-code-toggle">
              <input type="checkbox" bind:checked={includeSourceCode} />
              <span class="source-code-label">{labels.sourceCodeLabel}</span>
              <span class="source-code-note">{labels.sourceCodeNote}</span>
            </label>

            <hr class="divider" />

            <h3 class="plans-title">{labels.choosePlan}</h3>
            <div class="plans-grid">
              {#each plans as plan, i}
                <button
                  class="plan-card"
                  class:selected={selectedPlanKey === plan.key}
                  onclick={() => (selectedPlanKey = plan.key)}
                  type="button"
                  aria-pressed={selectedPlanKey === plan.key}
                  style="animation-delay: {i * 40}ms"
                >
                  {#if plan.badge}
                    <span class="plan-badge">{plan.badge[lang]}</span>
                  {/if}
                  <span class="plan-icon">{@html icons[plan.icon] ?? ''}</span>
                  <span class="plan-label">{plan.label[lang]}</span>
                  <span class="plan-primary">{plan.primary}</span>
                  {#if plan.secondary}
                    <span class="plan-secondary">{plan.secondary}</span>
                  {/if}
                  <span class="plan-desc">{plan.description[lang]}</span>
                  {#if selectedPlanKey === plan.key}
                    <span class="check-mark">
                      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                    </span>
                  {/if}
                </button>
              {/each}
            </div>

            <button class="includes-toggle" onclick={() => (showIncludes = !showIncludes)} type="button">
              {labels.includes} ({selectedSummary().length})
              <span class="includes-chevron" class:rotated={showIncludes}>{@html icons.chevronDown}</span>
            </button>
            {#if showIncludes}
              <ul class="includes-list">
                {#each selectedSummary() as item}
                  <li>
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>
                    {item}
                  </li>
                {/each}
              </ul>
            {/if}
          {/if}

          <hr class="divider" />

          <div class="contact-section">
            <h3 class="contact-title">{labels.contactTitle}</h3>

            <div class="form-grid">
              <div class="form-field" class:field-valid={fieldTouched.name && !formErrors.name && contactName.trim()} class:field-invalid={fieldTouched.name && formErrors.name}>
                <label for="q-name">{labels.name} *</label>
                <div class="input-wrapper">
                  <input id="q-name" type="text" bind:value={contactName} required aria-required="true" maxlength={200} onblur={() => { fieldTouched.name = true; validateField('name'); }} />
                  {#if fieldTouched.name && !formErrors.name && contactName.trim()}
                    <span class="field-icon field-icon-valid"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg></span>
                  {/if}
                </div>
                {#if fieldTouched.name && formErrors.name}<p class="field-error"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>{formErrors.name}</p>{/if}
              </div>
              <div class="form-field" class:field-valid={fieldTouched.email && !formErrors.email && contactEmail.trim()} class:field-invalid={fieldTouched.email && formErrors.email}>
                <label for="q-email">{labels.email} *</label>
                <div class="input-wrapper">
                  <input id="q-email" type="email" bind:value={contactEmail} required aria-required="true" maxlength={254} onblur={() => { fieldTouched.email = true; validateField('email'); }} />
                  {#if fieldTouched.email && !formErrors.email && contactEmail.trim()}
                    <span class="field-icon field-icon-valid"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg></span>
                  {/if}
                </div>
                {#if fieldTouched.email && formErrors.email}<p class="field-error"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>{formErrors.email}</p>{/if}
              </div>
              <div class="form-field">
                <label for="q-phone">{labels.phone} <span class="optional-hint">({lang === 'es' ? 'opcional' : 'optional'})</span></label>
                <input id="q-phone" type="tel" bind:value={contactPhone} maxlength={30} />
              </div>
              <div class="form-field">
                <label for="q-company">{labels.company} <span class="optional-hint">({lang === 'es' ? 'opcional' : 'optional'})</span></label>
                <input id="q-company" type="text" bind:value={contactCompany} maxlength={200} />
              </div>
              <div class="form-field full-width">
                <label for="q-notes">{labels.notes} <span class="optional-hint">({lang === 'es' ? 'opcional' : 'optional'})</span></label>
                <textarea id="q-notes" rows="3" bind:value={contactNotes} maxlength={2000}></textarea>
              </div>
              {#if turnstileSiteKey}
                <div class="full-width" use:initTurnstile></div>
              {/if}
            </div>

            {#if submitError}
              <p class="submit-error">{submitError}</p>
            {/if}

            <div class="submit-actions">
              <button class="quoter-submit-btn" onclick={submitQuote} disabled={submitting} type="button">
                {#if submitting}<span class="btn-spinner"></span>{/if}
                {submitting ? labels.sending : labels.send}
              </button>
            </div>
          </div>
        </div>
      {/if}
    </div>
    {/key}

    <!-- Navigation -->
    {#if currentStep < 7}
      <div class="step-nav">
        {#if currentStep > 1}
          <button class="quoter-nav-btn back" onclick={prevStep} type="button">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
            {labels.prev}
          </button>
        {:else}
          <div></div>
        {/if}
        <button class="quoter-nav-btn next" onclick={nextStep} disabled={!canNext()} type="button">
          {labels.next}
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
        </button>
      </div>
    {:else}
      <div class="step-nav">
        <button class="quoter-nav-btn back" onclick={prevStep} type="button">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="15 18 9 12 15 6"/></svg>
          {labels.prev}
        </button>
        <div></div>
      </div>
    {/if}
  </div>
{/if}

<style>
  /* ── Container ─────────────────────────────────────────── */
  .quoter-panel {
    background: var(--color-glass);
    backdrop-filter: blur(12px);
    border: 1px solid var(--color-glass-border);
    border-radius: 1.5rem;
    padding: 2rem;
    max-width: 820px;
    margin: 0 auto;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.06);
  }

  :global([data-theme="dark"]) .quoter-panel {
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  }

  /* ── Progress Bar ──────────────────────────────────────── */
  .progress-bar {
    height: 6px;
    background: var(--color-border);
    border-radius: 9999px;
    overflow: hidden;
    margin-bottom: 0.75rem;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--color-accent-primary), var(--color-accent-light));
    border-radius: 9999px;
    transition: width 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 0 8px rgba(37, 99, 235, 0.3);
  }

  .step-indicator {
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--color-text-muted);
    text-align: center;
    margin-bottom: 1.75rem;
  }

  /* ── Step Title / Hint ─────────────────────────────────── */
  .step-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.375rem;
    margin-bottom: 0.75rem;
    color: var(--color-text-primary);
    text-align: center;
  }

  .step-hint {
    font-size: 0.8125rem;
    color: var(--color-text-muted);
    margin-bottom: 1.25rem;
    text-align: center;
  }

  .step-content {
    min-height: 200px;
    margin-bottom: 1.5rem;
    animation: fadeSlideUp 0.3s ease-out;
  }

  /* ── Option Grid ───────────────────────────────────────── */
  .option-grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(3, 1fr);
  }

  .option-grid.cols-2 { grid-template-columns: repeat(2, 1fr); }
  .option-grid.cols-3 { grid-template-columns: repeat(3, 1fr); }

  @media (max-width: 767px) {
    .option-grid { grid-template-columns: repeat(2, 1fr); }
    .option-grid.cols-3 { grid-template-columns: repeat(2, 1fr); }
  }

  @media (max-width: 480px) {
    .option-grid { grid-template-columns: 1fr; }
    .option-grid.cols-2 { grid-template-columns: 1fr; }
    .option-grid.cols-3 { grid-template-columns: 1fr; }
  }

  /* ── Option Cards (Steps 1, 3-6) ──────────────────────── */
  .option-card {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1.25rem 1rem;
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    background: var(--color-bg-primary);
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-family: var(--font-sans);
    text-align: center;
    animation: fadeSlideUp 0.35s ease-out both;
  }

  .option-card:hover {
    border-color: var(--color-accent-light);
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(37, 99, 235, 0.08);
  }

  .option-card.selected {
    border-color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
    box-shadow: inset 0 0 0 1px var(--color-accent-primary);
  }

  /* ── Icon Containers ───────────────────────────────────── */
  .option-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 0.75rem;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .option-card.selected .option-icon {
    background: linear-gradient(135deg, var(--color-accent-primary), var(--color-accent-light));
    color: #fff;
  }

  .option-icon :global(svg) { width: 22px; height: 22px; }

  .option-label {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  .check-mark {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    color: var(--color-accent-primary);
  }

  /* ── Currency Card Extras ──────────────────────────────── */
  .currency-card { padding: 2rem 1.5rem; }
  .currency-flag { font-size: 2rem; }
  .currency-name { font-size: 0.75rem; color: var(--color-text-muted); }

  /* ── Step 2: Feature Groups & Cards ────────────────────── */
  .features-list {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .feature-group {
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    overflow: hidden;
  }

  .feature-group-header {
    display: flex;
    align-items: center;
    gap: 0.625rem;
    width: 100%;
    padding: 0.875rem 1rem;
    border: none;
    background: var(--color-bg-elevated);
    cursor: pointer;
    font-family: var(--font-sans);
    transition: background 0.2s;
  }

  .feature-group-header:hover {
    background: var(--color-accent-subtle);
  }

  .feature-group-icon {
    color: var(--color-accent-primary);
    flex-shrink: 0;
  }

  .feature-group-icon :global(svg) {
    width: 20px;
    height: 20px;
  }

  .feature-group-title {
    font-weight: 600;
    font-size: 0.9375rem;
    color: var(--color-text-primary);
    flex: 1;
    text-align: left;
  }

  .feature-group-count {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    flex-shrink: 0;
  }

  .feature-group-chevron {
    color: var(--color-text-muted);
    transition: transform 0.2s ease;
    flex-shrink: 0;
  }

  .feature-group-chevron :global(svg) {
    width: 18px;
    height: 18px;
  }

  .feature-group-chevron.collapsed {
    transform: rotate(-90deg);
  }

  .feature-cards-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
    padding: 1rem;
    background: transparent;
  }

  @media (max-width: 640px) {
    .feature-cards-grid {
      grid-template-columns: 1fr;
    }
  }

  .feature-card {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding: 0.875rem;
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    background: var(--color-glass);
    backdrop-filter: blur(8px);
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-family: var(--font-sans);
    text-align: left;
  }

  .feature-card:hover {
    border-color: var(--color-accent-light);
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(37, 99, 235, 0.08);
  }

  .feature-card.selected {
    border-color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
    box-shadow: inset 0 0 0 1px var(--color-accent-primary);
  }

  .feature-card-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.25rem;
  }

  .feature-card-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
    flex-shrink: 0;
  }

  .feature-card-icon :global(svg) {
    width: 16px;
    height: 16px;
  }

  .feature-badge {
    font-size: 0.6875rem;
    font-weight: 600;
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    flex-shrink: 0;
  }

  .feature-badge.included {
    color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
    text-transform: uppercase;
    letter-spacing: 0.02em;
  }

  .feature-badge.price {
    color: var(--color-text-muted);
    background: var(--color-border);
  }

  .feature-card-name {
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--color-text-primary);
    line-height: 1.3;
  }

  .feature-card-desc {
    font-size: 0.75rem;
    color: var(--color-text-muted);
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  /* ── Step 7: Result Section ────────────────────────────── */
  .estimate-display {
    position: relative;
    text-align: center;
    padding: 1.75rem 1.5rem;
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    background: var(--color-bg-primary);
    overflow: hidden;
  }

  .estimate-display::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--color-accent-primary), var(--color-accent-light));
  }

  .estimate-label {
    font-size: 0.8125rem;
    color: var(--color-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    margin-bottom: 0.5rem;
  }

  .estimate-range {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 700;
    color: var(--color-accent-primary);
  }

  /* ── Source Code Toggle ────────────────────────────────── */
  .source-code-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
    padding: 0.875rem 1rem;
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    background: var(--color-bg-primary);
    cursor: pointer;
    margin-top: 1rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .source-code-toggle:hover {
    border-color: var(--color-accent-light);
    background: var(--color-accent-subtle);
  }

  .source-code-toggle input[type="checkbox"] {
    accent-color: var(--color-accent-primary);
    width: 1rem;
    height: 1rem;
    cursor: pointer;
  }

  .source-code-label {
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  .source-code-note {
    font-size: 0.75rem;
    color: var(--color-text-muted);
  }

  /* ── Plans Section ─────────────────────────────────────── */
  .plans-title {
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 1rem;
    margin-bottom: 1rem;
    color: var(--color-text-primary);
    text-align: center;
  }

  .plans-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  @media (max-width: 640px) {
    .plans-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  .plan-card {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    padding: 1rem 0.75rem;
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    background: var(--color-bg-primary);
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-family: var(--font-sans);
    text-align: center;
    animation: fadeSlideUp 0.35s ease-out both;
  }

  .plan-card:hover {
    border-color: var(--color-accent-light);
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(37, 99, 235, 0.08);
  }

  .plan-card.selected {
    border-color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
    box-shadow: inset 0 0 0 1px var(--color-accent-primary);
  }

  .plan-badge {
    position: absolute;
    top: -0.5rem;
    right: 0.5rem;
    font-size: 0.625rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.03em;
    padding: 0.125rem 0.5rem;
    border-radius: 9999px;
    background: var(--color-accent-primary);
    color: #fff;
  }

  .plan-icon {
    color: var(--color-accent-primary);
  }

  .plan-icon :global(svg) { width: 20px; height: 20px; }

  .plan-label {
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--color-text-primary);
  }

  .plan-primary {
    font-family: var(--font-display);
    font-size: 1.125rem;
    font-weight: 700;
    color: var(--color-accent-primary);
    line-height: 1.2;
  }

  .plan-secondary {
    font-size: 0.6875rem;
    color: var(--color-text-muted);
  }

  .plan-desc {
    font-size: 0.6875rem;
    color: var(--color-text-muted);
    line-height: 1.3;
    margin-top: 0.125rem;
  }

  /* ── Collapsible Includes ──────────────────────────────── */
  .includes-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.75rem 0;
    border: none;
    background: transparent;
    cursor: pointer;
    font-family: var(--font-sans);
    font-size: 0.875rem;
    font-weight: 600;
    color: var(--color-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    transition: color 0.2s;
  }

  .includes-toggle:hover {
    color: var(--color-text-primary);
  }

  .includes-chevron {
    transition: transform 0.2s ease;
  }

  .includes-chevron :global(svg) { width: 16px; height: 16px; }

  .includes-chevron.rotated {
    transform: rotate(180deg);
  }

  .includes-list {
    list-style: none;
    padding: 0;
    display: grid;
    gap: 0.375rem;
    margin-bottom: 0.5rem;
  }

  .includes-list li {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.875rem;
    color: var(--color-text-secondary);
  }

  .includes-list li :global(svg) { color: var(--color-accent-primary); flex-shrink: 0; }

  .divider {
    border: none;
    border-top: 1px solid var(--color-border);
    margin: 2rem 0;
  }

  /* ── Contact Section ───────────────────────────────────── */
  .contact-section {
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    background: var(--color-bg-primary);
    padding: 1.5rem;
    animation: fadeSlideUp 0.35s ease-out both;
  }

  .contact-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.125rem;
    margin-bottom: 1.25rem;
    color: var(--color-text-primary);
    text-align: center;
  }

  .form-grid {
    display: grid;
    gap: 1rem;
    grid-template-columns: repeat(2, 1fr);
  }

  @media (max-width: 640px) {
    .form-grid { grid-template-columns: 1fr; }
  }

  .form-field { display: flex; flex-direction: column; gap: 0.25rem; }
  .form-field.full-width { grid-column: 1 / -1; }

  .form-field label {
    font-size: 0.8125rem;
    font-weight: 500;
    color: var(--color-text-secondary);
  }

  .optional-hint {
    font-weight: 400;
    font-style: italic;
    color: var(--color-text-muted);
  }

  .form-field input,
  .form-field textarea {
    padding: 0.75rem 1rem;
    border: 1.5px solid var(--color-border);
    border-radius: 0.625rem;
    background: var(--color-bg-elevated);
    color: var(--color-text-primary);
    font-family: var(--font-sans);
    font-size: 0.9375rem;
    min-height: 44px;
    transition: border-color 0.2s, box-shadow 0.2s;
  }

  .form-field input:focus,
  .form-field textarea:focus {
    outline: none;
    border-color: var(--color-accent-primary);
    box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
  }

  .input-wrapper {
    position: relative;
  }

  .input-wrapper input {
    width: 100%;
  }

  .field-icon {
    position: absolute;
    right: 0.75rem;
    top: 50%;
    transform: translateY(-50%);
    display: flex;
    pointer-events: none;
  }
  .field-icon-valid { color: var(--color-success); }

  .field-valid input { border-color: var(--color-success); }
  .field-invalid input { border-color: var(--color-error); }

  .field-error {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    font-size: 0.8125rem;
    color: var(--color-error);
    margin-top: 0.25rem;
  }

  .submit-error {
    font-size: 0.875rem;
    color: var(--color-error);
    text-align: center;
    margin: 1rem 0;
  }

  .submit-actions {
    display: flex;
    justify-content: center;
    margin-top: 1.5rem;
  }

  /* ── Submit Button ─────────────────────────────────────── */
  .quoter-submit-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.875rem 2rem;
    border: none;
    border-radius: 0.75rem;
    background: var(--color-accent-primary);
    color: #fff;
    font-family: var(--font-sans);
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
  }

  .quoter-submit-btn:hover:not(:disabled) {
    background: var(--color-accent-hover);
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(37, 99, 235, 0.3);
  }

  .quoter-submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  /* ── Navigation Buttons ────────────────────────────────── */
  .step-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 1rem;
    border-top: 1px solid var(--color-border);
  }

  .quoter-nav-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    padding: 0.625rem 1.25rem;
    border-radius: 0.75rem;
    font-family: var(--font-sans);
    font-size: 0.9375rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .quoter-nav-btn.back {
    background: transparent;
    border: 1.5px solid var(--color-border);
    color: var(--color-text-secondary);
  }

  .quoter-nav-btn.back:hover {
    border-color: var(--color-accent-light);
    color: var(--color-accent-primary);
    transform: translateX(-2px);
  }

  .quoter-nav-btn.next {
    background: var(--color-accent-primary);
    border: none;
    color: #fff;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.25);
  }

  .quoter-nav-btn.next:hover:not(:disabled) {
    background: var(--color-accent-hover);
    transform: translateY(-1px);
    box-shadow: 0 4px 16px rgba(37, 99, 235, 0.3);
  }

  .quoter-nav-btn.next:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .quoter-nav-btn :global(svg) {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
  }

  /* ── Success Screen ────────────────────────────────────── */
  .success-screen {
    position: relative;
    text-align: center;
    padding: 3.5rem 1.5rem;
  }

  .success-pulse {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 200px;
    height: 200px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(37, 99, 235, 0.08) 0%, transparent 70%);
    animation: successPulse 2s ease-in-out infinite;
    pointer-events: none;
  }

  .success-check {
    position: relative;
    color: var(--color-success);
    margin-bottom: 1.5rem;
    animation: popIn 0.4s ease-out;
  }

  .success-title {
    position: relative;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.75rem;
    margin-bottom: 0.75rem;
  }

  .success-msg {
    position: relative;
    color: var(--color-text-secondary);
    margin-bottom: 2rem;
  }

  .success-actions {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .success-actions .quoter-nav-btn {
    text-decoration: none;
  }

  /* ── Animations ────────────────────────────────────────── */
  @keyframes fadeSlideUp {
    from {
      opacity: 0;
      transform: translateY(8px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @keyframes popIn {
    0% { transform: scale(0); opacity: 0; }
    70% { transform: scale(1.1); }
    100% { transform: scale(1); opacity: 1; }
  }

  @keyframes successPulse {
    0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.5; }
    50% { transform: translate(-50%, -50%) scale(1.3); opacity: 0.2; }
  }

  @keyframes detailFadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @media (prefers-reduced-motion: reduce) {
    .step-content,
    .option-card,
    .plan-card,
    .contact-section,
    .success-pulse,
    .success-check,
    .progress-fill {
      animation: none !important;
      transition: none !important;
    }

    .option-card:hover,
    .feature-card:hover,
    .plan-card:hover,
    .quoter-nav-btn:hover {
      transform: none !important;
    }
  }

  /* ── Step 1: Floating Detail Tooltip (desktop only) ──────── */
  .option-card.expanded {
    z-index: 10;
  }

  .option-detail {
    position: absolute;
    top: calc(100% + 0.5rem);
    left: 0;
    width: 280px;
    z-index: 20;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 0.75rem;
    background: var(--color-bg-primary);
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
    animation: detailFadeIn 0.15s ease-out;
    pointer-events: none;
  }

  :global([data-theme="dark"]) .option-detail {
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }

  /* Right-align tooltip for cards in the last column */
  .option-grid > :nth-child(3n) .option-detail {
    left: auto;
    right: 0;
  }

  .option-detail-img {
    width: 100%;
    height: 120px;
    object-fit: cover;
    border-radius: 0.5rem;
  }

  .option-detail-desc {
    font-size: 0.8125rem;
    font-weight: 400;
    line-height: 1.5;
    color: var(--color-text-secondary);
    text-align: left;
  }

  /* Hide detail tooltip on mobile */
  @media (max-width: 767px) {
    .option-detail {
      display: none;
    }
  }

  /* ── Mobile Responsive ─────────────────────────────────── */
  @media (max-width: 767px) {
    .quoter-panel {
      padding: 1.25rem;
    }

    .step-title {
      font-size: 1.125rem;
    }

    .estimate-range {
      font-size: 1.5rem;
    }

    .contact-section {
      padding: 1.25rem;
    }
  }

  @media (max-width: 480px) {
    .plans-grid {
      grid-template-columns: 1fr;
    }

    .step-nav {
      flex-direction: column-reverse;
      gap: 0.5rem;
    }

    .quoter-nav-btn {
      width: 100%;
      justify-content: center;
    }
  }
</style>
