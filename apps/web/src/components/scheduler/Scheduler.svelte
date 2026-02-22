<script lang="ts">
  import { toast } from '../../lib/toast.svelte';

  interface Props {
    lang: 'es' | 'en';
    apiUrl?: string;
  }

  let { lang, apiUrl = '' }: Props = $props();

  const TOTAL_STEPS = 5;
  let currentStep = $state(1);

  // Step 1: Meeting type
  let meetingType = $state<'presencial' | 'videollamada' | ''>('');

  // Step 2-3: Calendar + time
  let today = new Date();
  let viewYear = $state(today.getFullYear());
  let viewMonth = $state(today.getMonth()); // 0-indexed
  let selectedDate = $state('');
  let selectedStartTime = $state('');
  let selectedEndTime = $state('');

  // Timezone
  let detectedTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  let clientTimezone = $state(detectedTimezone);
  let showTzPicker = $state(false);

  // Slots data
  let slots = $state<Array<{ date: string; startTime: string; endTime: string }>>([]);
  let loadingSlots = $state(false);
  let slotsError = $state('');

  // Step 4: Contact form
  let clientName = $state('');
  let clientEmail = $state('');
  let clientPhone = $state('');
  let clientCompany = $state('');
  let clientAddress = $state('');
  let notes = $state('');
  let formErrors = $state<Record<string, string>>({});
  let fieldTouched = $state<Record<string, boolean>>({});

  function validateField(field: string) {
    const errors = { ...formErrors };
    if (field === 'name') {
      if (!clientName.trim()) {
        errors.name = t.required;
      } else {
        delete errors.name;
      }
    }
    if (field === 'email') {
      if (!clientEmail.trim()) {
        errors.email = t.invalidEmail;
      } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(clientEmail.trim())) {
        errors.email = t.invalidEmail;
      } else {
        delete errors.email;
      }
    }
    formErrors = errors;
  }

  // Step 5: Confirmation
  let submitting = $state(false);
  let submitted = $state(false);
  let submitError = $state('');
  let bookingId = $state('');

  // i18n helper
  const t: Record<string, any> = lang === 'es' ? {
    title: 'Agenda una reunión',
    subtitle: 'Selecciona el tipo de reunión, fecha y horario que mejor te funcione.',
    step1Title: 'Tipo de reunión',
    inPerson: 'Presencial',
    inPersonDesc: 'Paso por tus oficinas',
    videoCall: 'Videollamada',
    videoCallDesc: 'Nos conectamos por video',
    step2Title: 'Selecciona una fecha',
    step3Title: 'Selecciona un horario',
    step4Title: 'Tus datos de contacto',
    name: 'Nombre',
    email: 'Email',
    phone: 'Teléfono',
    company: 'Empresa',
    address: 'Dirección',
    addressHint: 'Donde nos reuniremos',
    notes: 'Notas',
    notesPlaceholder: '¿Qué te gustaría discutir?',
    confirm: 'Enviar solicitud',
    step5Title: '¡Solicitud enviada!',
    confirmationMsg: 'Tu solicitud de reunión ha sido recibida. Te notificaré por email a',
    pendingNote: 'cuando sea confirmada.',
    bookingIdLabel: 'ID de reservación',
    backToHome: 'Volver al inicio',
    loading: 'Cargando disponibilidad...',
    noSlots: 'No hay horarios disponibles para este día.',
    noSlotsMonth: 'No hay disponibilidad este mes. Prueba el siguiente.',
    apiError: 'No se pudo conectar al servidor. Verifica que la API esté ejecutándose.',
    slotTaken: 'Este horario ya no está disponible. Por favor selecciona otro.',
    activeBooking: 'Ya tienes una solicitud activa. Espera a que sea procesada antes de agendar otra.',
    errorGeneric: 'Ocurrió un error. Inténtalo de nuevo.',
    next: 'Siguiente',
    back: 'Atrás',
    monthNames: ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio', 'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'],
    dayNames: ['Lu', 'Ma', 'Mi', 'Ju', 'Vi', 'Sa', 'Do'],
    meetingDate: 'Fecha',
    meetingTime: 'Hora',
    meetingType: 'Tipo',
    timezone: 'Zona horaria',
    changeTz: 'Cambiar',
    required: 'Este campo es requerido',
    invalidEmail: 'Email no válido',
  } : {
    title: 'Schedule a meeting',
    subtitle: 'Select the meeting type, date and time that works best for you.',
    step1Title: 'Meeting type',
    inPerson: 'In-person',
    inPersonDesc: "I'll visit your office",
    videoCall: 'Video call',
    videoCallDesc: 'We connect via video',
    step2Title: 'Select a date',
    step3Title: 'Select a time',
    step4Title: 'Your contact info',
    name: 'Name',
    email: 'Email',
    phone: 'Phone',
    company: 'Company',
    address: 'Address',
    addressHint: "Where we'll meet",
    notes: 'Notes',
    notesPlaceholder: 'What would you like to discuss?',
    confirm: 'Send request',
    step5Title: 'Request sent!',
    confirmationMsg: "Your meeting request has been received. I'll notify you by email at",
    pendingNote: 'when it\'s confirmed.',
    bookingIdLabel: 'Booking ID',
    backToHome: 'Back to home',
    loading: 'Loading availability...',
    noSlots: 'No time slots available for this day.',
    noSlotsMonth: 'No availability this month. Try the next one.',
    apiError: 'Could not connect to the server. Make sure the API is running.',
    slotTaken: 'This time slot is no longer available. Please select another.',
    activeBooking: 'You already have an active request. Wait for it to be processed before scheduling another.',
    errorGeneric: 'An error occurred. Please try again.',
    next: 'Next',
    back: 'Back',
    monthNames: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
    dayNames: ['Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa', 'Su'],
    meetingDate: 'Date',
    meetingTime: 'Time',
    meetingType: 'Type',
    timezone: 'Timezone',
    changeTz: 'Change',
    required: 'This field is required',
    invalidEmail: 'Invalid email',
  };

  // Derived: dates with available slots
  let availableDates = $derived(new Set(slots.map(s => s.date)));

  // Derived: slots for selected date
  let slotsForDate = $derived(slots.filter(s => s.date === selectedDate));

  // Calendar
  function getDaysInMonth(year: number, month: number) {
    return new Date(year, month + 1, 0).getDate();
  }

  function getFirstDayOfWeek(year: number, month: number) {
    const d = new Date(year, month, 1).getDay();
    return d === 0 ? 6 : d - 1; // Monday = 0
  }

  function isToday(year: number, month: number, day: number) {
    return year === today.getFullYear() && month === today.getMonth() && day === today.getDate();
  }

  function isPast(year: number, month: number, day: number) {
    const d = new Date(year, month, day);
    const todayStart = new Date(today.getFullYear(), today.getMonth(), today.getDate());
    return d < todayStart;
  }

  function dateStr(year: number, month: number, day: number) {
    return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
  }

  function formatDateDisplay(date: string) {
    const [y, m, d] = date.split('-');
    const monthName = t.monthNames[parseInt(m) - 1];
    const day = parseInt(d);
    if (lang === 'en') return `${monthName} ${day}, ${y}`;
    return `${day} de ${monthName.toLowerCase()}, ${y}`;
  }

  function formatTzOffset(tz: string) {
    try {
      const now = new Date();
      const formatter = new Intl.DateTimeFormat('en', { timeZone: tz, timeZoneName: 'shortOffset' });
      const parts = formatter.formatToParts(now);
      const offset = parts.find(p => p.type === 'timeZoneName')?.value || '';
      return offset; // e.g. "GMT-6"
    } catch {
      return '';
    }
  }

  // Common timezones for the Americas + worldwide
  const TIMEZONES = [
    'America/Mexico_City', 'America/Tijuana', 'America/Cancun', 'America/Hermosillo',
    'America/Monterrey', 'America/Chihuahua', 'America/Mazatlan',
    'America/New_York', 'America/Chicago', 'America/Denver', 'America/Los_Angeles',
    'America/Bogota', 'America/Lima', 'America/Santiago', 'America/Buenos_Aires',
    'America/Sao_Paulo', 'America/Toronto', 'America/Vancouver',
    'Europe/London', 'Europe/Madrid', 'Europe/Berlin', 'Europe/Paris',
    'Asia/Tokyo', 'Asia/Shanghai', 'Asia/Kolkata', 'Asia/Dubai',
    'Australia/Sydney', 'Pacific/Auckland',
  ];

  // Ensure detected timezone is in the list
  let tzOptions = $derived((() => {
    const set = new Set(TIMEZONES);
    if (!set.has(detectedTimezone)) set.add(detectedTimezone);
    return [...set].sort();
  })());

  function formatTime(time: string) {
    const [h, m] = time.split(':');
    const hour = parseInt(h);
    if (hour === 0) return `12:${m} AM`;
    if (hour < 12) return `${hour}:${m} AM`;
    if (hour === 12) return `12:${m} PM`;
    return `${hour - 12}:${m} PM`;
  }

  // Convert server time (America/Tijuana) to client timezone for display
  const SERVER_TZ = 'America/Tijuana';

  function toClientTime(date: string, time: string): string {
    if (clientTimezone === SERVER_TZ) return formatTime(time);
    try {
      const asUtc = new Date(`${date}T${time}:00Z`);
      const serverLocal = new Date(asUtc.toLocaleString('en-US', { timeZone: SERVER_TZ }));
      const offsetMs = asUtc.getTime() - serverLocal.getTime();
      const utcMoment = new Date(asUtc.getTime() + offsetMs);
      return utcMoment.toLocaleTimeString(lang === 'es' ? 'es-MX' : 'en-US', {
        timeZone: clientTimezone,
        hour: 'numeric',
        minute: '2-digit',
        hour12: true,
      });
    } catch {
      return formatTime(time);
    }
  }

  // Re-fetch slots when timezone changes
  let tzInitialized = false;
  $effect(() => {
    const _tz = clientTimezone;
    if (!tzInitialized) {
      tzInitialized = true;
      return;
    }
    if (currentStep >= 2) {
      selectedStartTime = '';
      selectedEndTime = '';
      fetchSlots();
    }
  });

  // Fetch slots for current view month
  async function fetchSlots() {
    loadingSlots = true;
    slotsError = '';
    const from = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-01`;
    const lastDay = getDaysInMonth(viewYear, viewMonth);
    const to = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-${String(lastDay).padStart(2, '0')}`;

    try {
      const res = await fetch(`${apiUrl}/scheduler/slots?from=${from}&to=${to}`);
      if (res.ok) {
        const data = await res.json();
        slots = data.slots || [];
      } else {
        slots = [];
        slotsError = t.apiError;
        toast.error(slotsError);
      }
    } catch {
      slots = [];
      slotsError = t.apiError;
      toast.error(slotsError);
    }
    loadingSlots = false;
  }

  function prevMonth() {
    if (viewMonth === today.getMonth() && viewYear === today.getFullYear()) return;
    if (viewMonth === 0) {
      viewMonth = 11;
      viewYear--;
    } else {
      viewMonth--;
    }
    selectedDate = '';
    selectedStartTime = '';
    selectedEndTime = '';
    fetchSlots();
  }

  function nextMonth() {
    if (viewMonth === 11) {
      viewMonth = 0;
      viewYear++;
    } else {
      viewMonth++;
    }
    selectedDate = '';
    selectedStartTime = '';
    selectedEndTime = '';
    fetchSlots();
  }

  function selectDate(day: number) {
    const ds = dateStr(viewYear, viewMonth, day);
    if (isPast(viewYear, viewMonth, day) || !availableDates.has(ds)) return;
    selectedDate = ds;
    selectedStartTime = '';
    selectedEndTime = '';
  }

  function selectSlot(slot: { startTime: string; endTime: string }) {
    selectedStartTime = slot.startTime;
    selectedEndTime = slot.endTime;
  }

  // Navigation
  let canNext = $derived((() => {
    switch (currentStep) {
      case 1: return meetingType !== '';
      case 2: return selectedDate !== '';
      case 3: return selectedStartTime !== '';
      case 4: return clientName.trim() !== '' && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(clientEmail.trim());
      default: return false;
    }
  })());

  function goNext() {
    if (!canNext) return;
    if (currentStep === 4) {
      submitBooking();
      return;
    }
    currentStep++;
    if (currentStep === 2) fetchSlots();
  }

  function goBack() {
    if (currentStep <= 1) return;
    submitError = '';
    if (currentStep === 3) {
      selectedStartTime = '';
      selectedEndTime = '';
      currentStep = 2;
      return;
    }
    currentStep--;
  }

  async function submitBooking() {
    formErrors = {};
    if (!clientName.trim()) formErrors.name = t.required;
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(clientEmail.trim())) formErrors.email = t.invalidEmail;
    if (Object.keys(formErrors).length > 0) return;

    submitting = true;
    submitError = '';

    try {
      const res = await fetch(`${apiUrl}/scheduler/bookings`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          date: selectedDate,
          startTime: selectedStartTime,
          meetingType,
          clientName: clientName.trim(),
          clientEmail: clientEmail.trim(),
          clientPhone: clientPhone.trim(),
          clientCompany: clientCompany.trim(),
          clientAddress: clientAddress.trim(),
          clientTimezone: clientTimezone,
          notes: notes.trim(),
          lang,
        }),
      });

      const data = await res.json();

      if (res.status === 409) {
        // Could be slot taken OR active booking exists
        submitError = data.message || t.slotTaken;
        toast.error(submitError);
        submitting = false;
        // If slot taken, go back to time selection
        if (!data.message?.includes('activa') && !data.message?.includes('active')) {
          currentStep = 2;
          selectedStartTime = '';
          selectedEndTime = '';
          fetchSlots();
        }
        return;
      }

      if (!res.ok || !data.success) {
        submitError = data.message || t.errorGeneric;
        toast.error(submitError);
        submitting = false;
        return;
      }

      bookingId = data.bookingId;
      submitted = true;
      currentStep = 5;
      toast.success(t.step5Title);
    } catch {
      submitError = t.errorGeneric;
      toast.error(submitError);
    }
    submitting = false;
  }
</script>

<div class="scheduler">
  <!-- Header -->
  <div class="scheduler-header">
    <h2 class="scheduler-title">{t.title}</h2>
    <p class="scheduler-subtitle">{t.subtitle}</p>
  </div>

  <!-- Progress -->
  <div class="progress">
    {#each Array(TOTAL_STEPS) as _, i}
      <div class="progress-step" class:active={currentStep >= i + 1} class:current={currentStep === i + 1}>
        <span>{i + 1}</span>
      </div>
      {#if i < TOTAL_STEPS - 1}
        <div class="progress-line" class:active={currentStep > i + 1}></div>
      {/if}
    {/each}
  </div>

  <!-- Steps -->
  <div class="step-container">
    <!-- STEP 1: Meeting Type -->
    {#if currentStep === 1}
      <div class="step" style="animation: slideIn 0.3s ease-out">
        <h3 class="step-title">{t.step1Title}</h3>
        <div class="type-grid">
          <button
            type="button"
            class="type-card"
            class:selected={meetingType === 'presencial'}
            onclick={() => { meetingType = 'presencial'; goNext(); }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
            <span class="type-label">{t.inPerson}</span>
            <span class="type-desc">{t.inPersonDesc}</span>
          </button>
          <button
            type="button"
            class="type-card"
            class:selected={meetingType === 'videollamada'}
            onclick={() => { meetingType = 'videollamada'; goNext(); }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>
            <span class="type-label">{t.videoCall}</span>
            <span class="type-desc">{t.videoCallDesc}</span>
          </button>
        </div>
      </div>
    {/if}

    <!-- STEP 2: Calendar -->
    {#if currentStep === 2}
      <div class="step" style="animation: slideIn 0.3s ease-out">
        <h3 class="step-title">{t.step2Title}</h3>
        <div class="calendar">
          <div class="calendar-nav">
            <button type="button" class="cal-nav-btn" onclick={prevMonth}
              disabled={viewMonth === today.getMonth() && viewYear === today.getFullYear()}>
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
            </button>
            <span class="cal-month">{t.monthNames[viewMonth]} {viewYear}</span>
            <button type="button" class="cal-nav-btn" onclick={nextMonth}>
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
            </button>
          </div>

          {#if loadingSlots}
            <div class="loading">{t.loading}</div>
          {:else}
            <div class="calendar-grid">
              {#each t.dayNames as dayName}
                <div class="cal-day-name">{dayName}</div>
              {/each}

              {#each Array(getFirstDayOfWeek(viewYear, viewMonth)) as _}
                <div class="cal-day empty"></div>
              {/each}

              {#each Array(getDaysInMonth(viewYear, viewMonth)) as _, i}
                {@const day = i + 1}
                {@const ds = dateStr(viewYear, viewMonth, day)}
                {@const past = isPast(viewYear, viewMonth, day)}
                {@const hasSlots = availableDates.has(ds)}
                <button
                  type="button"
                  class="cal-day"
                  class:past
                  class:available={hasSlots && !past}
                  class:selected={selectedDate === ds}
                  class:today={isToday(viewYear, viewMonth, day)}
                  disabled={past || !hasSlots}
                  onclick={() => selectDate(day)}
                >
                  {day}
                </button>
              {/each}
            </div>

            {#if !loadingSlots && slotsError}
              <div class="error-msg">{slotsError}</div>
            {:else if !loadingSlots && slots.length === 0}
              <p class="no-slots-msg">{t.noSlotsMonth}</p>
            {/if}
          {/if}
        </div>

        <!-- Timezone indicator -->
        <div class="tz-bar">
          <span class="tz-label">{t.timezone}:</span>
          {#if showTzPicker}
            <select class="tz-select" bind:value={clientTimezone} onchange={() => showTzPicker = false}>
              {#each tzOptions as tz}
                <option value={tz}>{tz.replace(/_/g, ' ')} ({formatTzOffset(tz)})</option>
              {/each}
            </select>
          {:else}
            <span class="tz-value">{clientTimezone.replace(/_/g, ' ')} ({formatTzOffset(clientTimezone)})</span>
            <button type="button" class="tz-change" onclick={() => showTzPicker = true}>{t.changeTz}</button>
          {/if}
        </div>

        <!-- Time slots shown below calendar when date is selected -->
        {#if selectedDate && slotsForDate.length > 0}
          <div class="time-section" style="animation: slideIn 0.3s ease-out">
            <h3 class="step-title">{t.step3Title}</h3>
            <div class="time-grid">
              {#each slotsForDate as slot}
                <button
                  type="button"
                  class="time-btn"
                  class:selected={selectedStartTime === slot.startTime}
                  onclick={() => { selectSlot(slot); currentStep = 3; }}
                >
                  {toClientTime(slot.date, slot.startTime)}
                </button>
              {/each}
            </div>
          </div>
        {:else if selectedDate && slotsForDate.length === 0}
          <p class="no-slots-msg">{t.noSlots}</p>
        {/if}
      </div>
    {/if}

    <!-- STEP 3: Time confirmation (brief, then auto-advance to step 4) -->
    {#if currentStep === 3}
      <div class="step" style="animation: slideIn 0.3s ease-out">
        <h3 class="step-title">{t.step3Title}</h3>
        <div class="selected-summary">
          <div class="summary-item">
            <span class="summary-label">{t.meetingDate}</span>
            <span class="summary-value">{formatDateDisplay(selectedDate)}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">{t.meetingTime}</span>
            <span class="summary-value">{toClientTime(selectedDate, selectedStartTime)} - {toClientTime(selectedDate, selectedEndTime)}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">{t.timezone}</span>
            <span class="summary-value">{clientTimezone.replace(/_/g, ' ')} ({formatTzOffset(clientTimezone)})</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">{t.meetingType}</span>
            <span class="summary-value">{meetingType === 'presencial' ? t.inPerson : t.videoCall}</span>
          </div>
        </div>
      </div>
    {/if}

    <!-- STEP 4: Contact Form -->
    {#if currentStep === 4}
      <div class="step" style="animation: slideIn 0.3s ease-out">
        <h3 class="step-title">{t.step4Title}</h3>
        <form class="contact-form" onsubmit={(e) => { e.preventDefault(); submitBooking(); }}>
          <div class="form-field" class:field-valid={fieldTouched.name && !formErrors.name && clientName.trim()} class:field-invalid={fieldTouched.name && formErrors.name}>
            <label for="s-name">{t.name} *</label>
            <div class="input-wrapper">
              <input id="s-name" type="text" bind:value={clientName} required autocomplete="name" onblur={() => { fieldTouched.name = true; validateField('name'); }} />
              {#if fieldTouched.name && !formErrors.name && clientName.trim()}
                <span class="field-icon field-icon-valid"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg></span>
              {/if}
            </div>
            {#if fieldTouched.name && formErrors.name}<span class="field-error"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>{formErrors.name}</span>{/if}
          </div>
          <div class="form-field" class:field-valid={fieldTouched.email && !formErrors.email && clientEmail.trim()} class:field-invalid={fieldTouched.email && formErrors.email}>
            <label for="s-email">{t.email} *</label>
            <div class="input-wrapper">
              <input id="s-email" type="email" bind:value={clientEmail} required autocomplete="email" onblur={() => { fieldTouched.email = true; validateField('email'); }} />
              {#if fieldTouched.email && !formErrors.email && clientEmail.trim()}
                <span class="field-icon field-icon-valid"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg></span>
              {/if}
            </div>
            {#if fieldTouched.email && formErrors.email}<span class="field-error"><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>{formErrors.email}</span>{/if}
          </div>
          <div class="form-field">
            <label for="s-phone">{t.phone}</label>
            <input id="s-phone" type="tel" bind:value={clientPhone} autocomplete="tel" />
          </div>
          <div class="form-field">
            <label for="s-company">{t.company}</label>
            <input id="s-company" type="text" bind:value={clientCompany} autocomplete="organization" />
          </div>
          {#if meetingType === 'presencial'}
            <div class="form-field">
              <label for="s-address">{t.address}</label>
              <input id="s-address" type="text" bind:value={clientAddress} placeholder={t.addressHint} autocomplete="street-address" />
            </div>
          {/if}
          <div class="form-field">
            <label for="s-notes">{t.notes}</label>
            <textarea id="s-notes" bind:value={notes} placeholder={t.notesPlaceholder} rows="3"></textarea>
          </div>
          {#if submitError}
            <div class="error-msg">{submitError}</div>
          {/if}
          <button type="submit" class="btn-primary" style="width: 100%;" disabled={submitting}>
            {#if submitting}
              <span class="spinner"></span>
            {:else}
              {t.confirm}
            {/if}
          </button>
        </form>
      </div>
    {/if}

    <!-- STEP 5: Pending Confirmation -->
    {#if currentStep === 5}
      <div class="step confirmation" style="animation: slideIn 0.3s ease-out">
        <div class="success-check">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
        </div>
        <h3 class="success-title">{t.step5Title}</h3>
        <div class="pending-badge">
          {lang === 'es' ? 'Pendiente de confirmación' : 'Pending confirmation'}
        </div>
        <div class="success-details">
          <div class="summary-item">
            <span class="summary-label">{t.meetingDate}</span>
            <span class="summary-value">{formatDateDisplay(selectedDate)}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">{t.meetingTime}</span>
            <span class="summary-value">{toClientTime(selectedDate, selectedStartTime)} - {toClientTime(selectedDate, selectedEndTime)}</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">{t.timezone}</span>
            <span class="summary-value">{clientTimezone.replace(/_/g, ' ')} ({formatTzOffset(clientTimezone)})</span>
          </div>
          <div class="summary-item">
            <span class="summary-label">{t.meetingType}</span>
            <span class="summary-value">{meetingType === 'presencial' ? t.inPerson : t.videoCall}</span>
          </div>
        </div>
        <p class="confirmation-email">{t.confirmationMsg} <strong>{clientEmail}</strong> {t.pendingNote}</p>
        <p class="booking-id">{t.bookingIdLabel}: <strong>{bookingId}</strong></p>
        <a href={lang === 'es' ? '/es/' : '/en/'} class="home-btn">{t.backToHome}</a>
      </div>
    {/if}
  </div>

  <!-- Navigation buttons -->
  {#if currentStep > 1 && currentStep < 5}
    <div class="nav-buttons">
      <button type="button" class="nav-btn back" onclick={goBack}>{t.back}</button>
      {#if currentStep === 3}
        <button type="button" class="nav-btn next" onclick={goNext} disabled={!canNext}>{t.next}</button>
      {/if}
    </div>
  {/if}
</div>

<style>
  .scheduler {
    background: var(--color-glass);
    backdrop-filter: blur(12px);
    border: 1px solid var(--color-glass-border);
    border-radius: 1.25rem;
    padding: 2rem;
    max-width: 600px;
    margin: 0 auto;
  }

  .scheduler-header {
    text-align: center;
    margin-bottom: 1.5rem;
  }

  .scheduler-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.75rem;
    margin-bottom: 0.5rem;
  }

  .scheduler-subtitle {
    color: var(--color-text-secondary);
    font-size: 0.9375rem;
  }

  /* Progress */
  .progress {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0;
    margin-bottom: 2rem;
  }

  .progress-step {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 600;
    border: 2px solid var(--color-border);
    color: var(--color-text-muted);
    background: var(--color-bg-primary);
    transition: all 0.3s;
  }

  .progress-step.active {
    border-color: var(--color-accent-primary);
    color: var(--color-accent-primary);
  }

  .progress-step.current {
    background: var(--color-accent-primary);
    border-color: var(--color-accent-primary);
    color: #fff;
  }

  .progress-line {
    width: 2rem;
    height: 2px;
    background: var(--color-border);
    transition: background 0.3s;
  }

  .progress-line.active {
    background: var(--color-accent-primary);
  }

  /* Steps */
  .step-container {
    min-height: 300px;
  }

  .step-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.25rem;
    margin-bottom: 1.25rem;
    text-align: center;
  }

  /* Step 1: Type selection */
  .type-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .type-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 1.5rem 1rem;
    background: var(--color-bg-primary);
    border: 2px solid var(--color-border);
    border-radius: 1rem;
    cursor: pointer;
    transition: all 0.2s;
    color: var(--color-text-primary);
  }

  .type-card:hover {
    border-color: var(--color-accent-primary);
  }

  .type-card.selected {
    border-color: var(--color-accent-primary);
    background: var(--color-accent-subtle);
  }

  .type-label {
    font-weight: 600;
    font-size: 1rem;
  }

  .type-desc {
    font-size: 0.8125rem;
    color: var(--color-text-secondary);
  }

  /* Calendar */
  .calendar {
    background: var(--color-bg-primary);
    border: 1px solid var(--color-border);
    border-radius: 1rem;
    padding: 1.25rem;
  }

  .calendar-nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
  }

  .cal-nav-btn {
    background: none;
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    padding: 0.625rem;
    min-width: 44px;
    min-height: 44px;
    cursor: pointer;
    color: var(--color-text-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
  }

  .cal-nav-btn:hover:not(:disabled) {
    border-color: var(--color-accent-primary);
    color: var(--color-accent-primary);
  }

  .cal-nav-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .cal-month {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.0625rem;
  }

  .calendar-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 0.25rem;
  }

  .cal-day-name {
    text-align: center;
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--color-text-muted);
    padding: 0.375rem 0;
  }

  .cal-day {
    aspect-ratio: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    border: 2px solid transparent;
    border-radius: 0.5rem;
    background: none;
    color: var(--color-text-primary);
    cursor: default;
    transition: all 0.15s;
  }

  .cal-day.empty {
    visibility: hidden;
  }

  .cal-day.past {
    color: var(--color-text-muted);
    opacity: 0.4;
  }

  .cal-day.available {
    cursor: pointer;
    background: var(--color-accent-subtle);
    color: var(--color-accent-primary);
    font-weight: 600;
  }

  .cal-day.available:hover {
    border-color: var(--color-accent-primary);
  }

  .cal-day.selected {
    background: var(--color-accent-primary);
    color: #fff;
    border-color: var(--color-accent-primary);
  }

  .cal-day.today {
    border-color: var(--color-accent-light);
  }

  .loading {
    text-align: center;
    color: var(--color-text-muted);
    padding: 2rem;
  }

  .no-slots-msg {
    text-align: center;
    color: var(--color-text-muted);
    font-size: 0.875rem;
    margin-top: 1rem;
  }

  /* Timezone bar */
  .tz-bar {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-top: 1rem;
    padding: 0.625rem 0.875rem;
    background: var(--color-accent-subtle);
    border-radius: 0.5rem;
    font-size: 0.8125rem;
    flex-wrap: wrap;
  }

  .tz-label {
    color: var(--color-text-secondary);
    font-weight: 600;
    white-space: nowrap;
  }

  .tz-value {
    color: var(--color-text-primary);
    font-weight: 500;
  }

  .tz-change {
    background: none;
    border: none;
    color: var(--color-accent-primary);
    font-size: 0.8125rem;
    font-weight: 600;
    cursor: pointer;
    text-decoration: underline;
    padding: 0;
  }

  .tz-select {
    flex: 1;
    min-width: 200px;
    padding: 0.375rem 0.5rem;
    background: var(--color-bg-primary);
    border: 1px solid var(--color-border);
    border-radius: 0.375rem;
    font-size: 0.8125rem;
    color: var(--color-text-primary);
    font-family: inherit;
  }

  .tz-select:focus {
    outline: none;
    border-color: var(--color-accent-primary);
  }

  /* Time slots */
  .time-section {
    margin-top: 1.5rem;
  }

  .time-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.5rem;
  }

  .time-btn {
    padding: 0.75rem;
    background: var(--color-bg-primary);
    border: 2px solid var(--color-border);
    border-radius: 0.75rem;
    cursor: pointer;
    font-weight: 500;
    font-size: 0.9375rem;
    color: var(--color-text-primary);
    transition: all 0.15s;
  }

  .time-btn:hover {
    border-color: var(--color-accent-primary);
    color: var(--color-accent-primary);
  }

  .time-btn.selected {
    background: var(--color-accent-primary);
    border-color: var(--color-accent-primary);
    color: #fff;
  }

  /* Summary */
  .selected-summary, .success-details {
    background: var(--color-bg-primary);
    border: 1px solid var(--color-border);
    border-radius: 1rem;
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .summary-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .summary-label {
    font-size: 0.875rem;
    color: var(--color-text-secondary);
  }

  .summary-value {
    font-weight: 600;
    font-size: 0.9375rem;
  }

  /* Contact form */
  .contact-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .form-field {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
  }

  .form-field label {
    font-size: 0.8125rem;
    font-weight: 600;
    color: var(--color-text-secondary);
  }

  .form-field input,
  .form-field textarea {
    padding: 0.75rem 1rem;
    background: var(--color-bg-elevated);
    border: 1.5px solid var(--color-border);
    border-radius: 0.5rem;
    font-size: 0.9375rem;
    min-height: 44px;
    color: var(--color-text-primary);
    font-family: var(--font-sans);
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

  .error-msg {
    background: color-mix(in srgb, var(--color-error) 10%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-error) 30%, transparent);
    border-radius: 0.5rem;
    padding: 0.75rem;
    color: var(--color-error);
    font-size: 0.875rem;
    text-align: center;
  }

  .spinner {
    width: 1.25rem;
    height: 1.25rem;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }

  /* Confirmation */
  .confirmation {
    text-align: center;
  }

  .success-check {
    color: var(--color-success);
    margin-bottom: 1rem;
    animation: popIn 0.4s ease-out;
    display: flex;
    justify-content: center;
  }

  .success-title {
    font-family: var(--font-display);
    font-weight: 700;
    font-size: 1.5rem;
    margin-bottom: 0.75rem;
  }

  .pending-badge {
    display: inline-block;
    padding: 0.375rem 1rem;
    background: color-mix(in srgb, var(--color-success) 15%, transparent);
    color: var(--color-success);
    border-radius: 2rem;
    font-size: 0.8125rem;
    font-weight: 600;
    margin-bottom: 1.25rem;
  }

  .success-details {
    margin-bottom: 1.25rem;
    text-align: left;
  }

  .confirmation-email {
    color: var(--color-text-secondary);
    font-size: 0.9375rem;
    margin-bottom: 0.5rem;
  }

  .booking-id {
    color: var(--color-text-muted);
    font-size: 0.875rem;
    margin-bottom: 1.5rem;
  }

  .home-btn {
    display: inline-block;
    padding: 0.75rem 1.5rem;
    background: var(--color-accent-primary);
    color: #fff;
    border-radius: 0.75rem;
    text-decoration: none;
    font-weight: 600;
    transition: opacity 0.2s;
  }

  .home-btn:hover {
    opacity: 0.9;
  }

  /* Navigation */
  .nav-buttons {
    display: flex;
    justify-content: space-between;
    margin-top: 1.5rem;
    gap: 1rem;
  }

  .nav-btn {
    padding: 0.625rem 1.25rem;
    border-radius: 0.5rem;
    font-weight: 500;
    font-size: 0.875rem;
    cursor: pointer;
    transition: all 0.2s;
  }

  .nav-btn.back {
    background: none;
    border: 1px solid var(--color-border);
    color: var(--color-text-secondary);
  }

  .nav-btn.back:hover {
    border-color: var(--color-text-primary);
    color: var(--color-text-primary);
  }

  .nav-btn.next {
    background: var(--color-accent-primary);
    border: none;
    color: #fff;
    margin-left: auto;
  }

  .nav-btn.next:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  @keyframes slideIn {
    from { opacity: 0; transform: translateX(20px); }
    to { opacity: 1; transform: translateX(0); }
  }

  @keyframes popIn {
    0% { transform: scale(0); opacity: 0; }
    70% { transform: scale(1.1); }
    100% { transform: scale(1); opacity: 1; }
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  /* Mobile */
  @media (max-width: 640px) {
    .scheduler {
      padding: 1.25rem;
      border-radius: 0.75rem;
    }

    .scheduler-title {
      font-size: 1.375rem;
    }

    .type-grid {
      grid-template-columns: 1fr;
    }

    .time-grid {
      grid-template-columns: repeat(2, 1fr);
    }

    .progress-step {
      width: 1.75rem;
      height: 1.75rem;
      font-size: 0.6875rem;
    }

    .progress-line {
      width: 1.25rem;
    }
  }
</style>
