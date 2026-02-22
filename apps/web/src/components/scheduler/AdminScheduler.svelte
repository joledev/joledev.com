<script lang="ts">
  interface Props {
    lang: 'es' | 'en';
    apiUrl?: string;
  }

  let { lang, apiUrl = '' }: Props = $props();

  let password = $state('');
  let authenticated = $state(false);
  let authError = $state('');

  // Calendar state
  let today = new Date();
  let viewYear = $state(today.getFullYear());
  let viewMonth = $state(today.getMonth());

  // Bookings data
  type BookingData = {
    id: number;
    bookingId: string;
    date: string;
    startTime: string;
    endTime: string;
    meetingType: string;
    clientName: string;
    clientEmail: string;
    clientPhone: string;
    clientCompany: string;
    clientAddress: string;
    clientTimezone: string;
    notes: string;
    lang: string;
    status: string;
    createdAt: string;
  };

  let bookings = $state<BookingData[]>([]);
  let loading = $state(false);
  let selectedBooking = $state<BookingData | null>(null);
  let statusMsg = $state('');

  const isEs = lang === 'es';
  const labels = {
    title: isEs ? 'Admin — Agenda' : 'Admin — Schedule',
    password: isEs ? 'Contraseña' : 'Password',
    login: isEs ? 'Entrar' : 'Login',
    cancelBooking: isEs ? 'Cancelar reservación' : 'Cancel booking',
    close: isEs ? 'Cerrar' : 'Close',
    pending: isEs ? 'Pendiente' : 'Pending',
    confirmed: isEs ? 'Confirmada' : 'Confirmed',
    rejected: isEs ? 'Rechazada' : 'Rejected',
    cancelled: isEs ? 'Cancelada' : 'Cancelled',
    wrongPassword: isEs ? 'Contraseña incorrecta' : 'Wrong password',
    noBookings: isEs ? 'Sin reservaciones este mes' : 'No bookings this month',
  };

  const statusLabels: Record<string, string> = {
    pending: labels.pending,
    confirmed: labels.confirmed,
    rejected: labels.rejected,
    cancelled: labels.cancelled,
  };

  const dayLabels = isEs
    ? ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']
    : ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];

  const monthNames = isEs
    ? ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio', 'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre']
    : ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];

  function authHeaders() {
    return { 'Authorization': 'Basic ' + btoa('admin:' + password), 'Content-Type': 'application/json' };
  }

  async function tryLogin() {
    authError = '';
    try {
      const from = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-01`;
      const to = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-28`;
      const res = await fetch(`${apiUrl}/scheduler/admin/bookings?from=${from}&to=${to}`, {
        headers: authHeaders(),
      });
      if (res.status === 401) {
        authError = labels.wrongPassword;
        return;
      }
      authenticated = true;
      sessionStorage.setItem('scheduler-admin-pw', password);
      fetchBookings();
    } catch {
      authError = 'Error connecting to API';
    }
  }

  async function fetchBookings() {
    loading = true;
    const daysInMonth = new Date(viewYear, viewMonth + 1, 0).getDate();
    const from = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-01`;
    const to = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-${String(daysInMonth).padStart(2, '0')}`;

    try {
      const res = await fetch(`${apiUrl}/scheduler/admin/bookings?from=${from}&to=${to}`, {
        headers: authHeaders(),
      });
      if (res.ok) {
        const data = await res.json();
        bookings = data.bookings || [];
      }
    } catch { /* ignore */ }
    loading = false;
  }

  function prevMonth() {
    if (viewMonth === 0) { viewMonth = 11; viewYear--; } else { viewMonth--; }
    fetchBookings();
  }

  function nextMonth() {
    if (viewMonth === 11) { viewMonth = 0; viewYear++; } else { viewMonth++; }
    fetchBookings();
  }

  function bookingsForDay(day: number): BookingData[] {
    const ds = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
    return bookings.filter(b => b.date === ds);
  }

  function getDaysInMonth(y: number, m: number) { return new Date(y, m + 1, 0).getDate(); }
  function getFirstDow(y: number, m: number) { const d = new Date(y, m, 1).getDay(); return d === 0 ? 6 : d - 1; }

  function formatTime(t: string) {
    const [h, m] = t.split(':');
    const hr = parseInt(h);
    if (hr < 12) return `${hr}:${m}`;
    if (hr === 12) return `12:${m}`;
    return `${hr - 12}:${m}p`;
  }

  async function cancelBooking(bookingDbId: number) {
    try {
      const res = await fetch(`${apiUrl}/scheduler/admin/bookings/${bookingDbId}`, {
        method: 'PATCH',
        headers: authHeaders(),
        body: JSON.stringify({ status: 'cancelled' }),
      });
      if (res.ok) {
        statusMsg = 'Booking cancelled';
        selectedBooking = null;
        fetchBookings();
      }
    } catch { /* ignore */ }
  }

  // Restore password from session
  $effect(() => {
    const saved = sessionStorage.getItem('scheduler-admin-pw');
    if (saved) {
      password = saved;
      tryLogin();
    }
  });
</script>

{#if !authenticated}
  <div class="admin-login">
    <h1>{labels.title}</h1>
    <form onsubmit={(e) => { e.preventDefault(); tryLogin(); }}>
      <input type="password" bind:value={password} placeholder={labels.password} autocomplete="current-password" />
      <button type="submit">{labels.login}</button>
    </form>
    {#if authError}<p class="error">{authError}</p>{/if}
  </div>
{:else}
  <div class="admin">
    <h1>{labels.title}</h1>

    {#if statusMsg}
      <div class="status-msg" onclick={() => statusMsg = ''} onkeydown={() => {}}>{statusMsg}</div>
    {/if}

    <!-- Legend -->
    <div class="legend">
      <span class="legend-item"><span class="dot pending"></span> {labels.pending}</span>
      <span class="legend-item"><span class="dot confirmed"></span> {labels.confirmed}</span>
      <span class="legend-item"><span class="dot rejected"></span> {labels.rejected}</span>
      <span class="legend-item"><span class="dot cancelled"></span> {labels.cancelled}</span>
    </div>

    <!-- Calendar View -->
    <div class="cal-header">
      <button type="button" onclick={prevMonth}>&lt;</button>
      <h2>{monthNames[viewMonth]} {viewYear}</h2>
      <button type="button" onclick={nextMonth}>&gt;</button>
    </div>

    {#if loading}
      <p class="loading">Loading...</p>
    {:else}
      <div class="admin-cal">
        {#each dayLabels as dName}
          <div class="cal-head">{dName}</div>
        {/each}
        {#each Array(getFirstDow(viewYear, viewMonth)) as _}
          <div class="cal-cell empty"></div>
        {/each}
        {#each Array(getDaysInMonth(viewYear, viewMonth)) as _, i}
          {@const day = i + 1}
          {@const dayBookings = bookingsForDay(day)}
          <div class="cal-cell" class:has-bookings={dayBookings.length > 0}>
            <div class="cell-day">{day}</div>
            {#each dayBookings as booking}
              <!-- svelte-ignore a11y_no_static_element_interactions -->
              <div
                class="booking-chip {booking.status}"
                onclick={() => selectedBooking = booking}
                onkeydown={() => {}}
              >
                <span class="chip-time">{formatTime(booking.startTime)}</span>
                <span class="chip-name">{booking.clientName.split(' ')[0]}</span>
              </div>
            {/each}
          </div>
        {/each}
      </div>

      {#if bookings.length === 0}
        <p class="no-bookings">{labels.noBookings}</p>
      {/if}
    {/if}

    <!-- Booking Detail Modal -->
    {#if selectedBooking}
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div class="modal-overlay" onclick={() => selectedBooking = null} onkeydown={() => {}}>
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="modal" onclick={(e) => e.stopPropagation()} onkeydown={() => {}}>
          <h3>{selectedBooking.bookingId}</h3>
          <p><strong>Client:</strong> {selectedBooking.clientName}</p>
          <p><strong>Email:</strong> {selectedBooking.clientEmail}</p>
          {#if selectedBooking.clientPhone}
            <p><strong>Phone:</strong> {selectedBooking.clientPhone}</p>
          {/if}
          {#if selectedBooking.clientCompany}
            <p><strong>Company:</strong> {selectedBooking.clientCompany}</p>
          {/if}
          <p><strong>Type:</strong> {selectedBooking.meetingType}</p>
          {#if selectedBooking.clientAddress}
            <p><strong>Address:</strong> {selectedBooking.clientAddress}</p>
          {/if}
          <p><strong>Date:</strong> {selectedBooking.date} {selectedBooking.startTime}-{selectedBooking.endTime}</p>
          {#if selectedBooking.clientTimezone}
            <p><strong>Timezone:</strong> {selectedBooking.clientTimezone}</p>
          {/if}
          {#if selectedBooking.notes}
            <p><strong>Notes:</strong> {selectedBooking.notes}</p>
          {/if}
          <p><strong>Status:</strong>
            <span class="status-label {selectedBooking.status}">
              {statusLabels[selectedBooking.status] || selectedBooking.status}
            </span>
          </p>
          <div class="modal-actions">
            {#if selectedBooking.status === 'pending' || selectedBooking.status === 'confirmed'}
              <button type="button" class="cancel-btn" onclick={() => selectedBooking && cancelBooking(selectedBooking.id)}>
                {labels.cancelBooking}
              </button>
            {/if}
            <button type="button" class="close-btn" onclick={() => selectedBooking = null}>{labels.close}</button>
          </div>
        </div>
      </div>
    {/if}
  </div>
{/if}

<style>
  .admin-login {
    max-width: 320px;
    margin: 4rem auto;
    text-align: center;
  }

  .admin-login h1 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 1.5rem;
  }

  .admin-login form {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .admin-login input {
    padding: 0.625rem;
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    font-size: 1rem;
    background: var(--color-bg-primary);
    color: var(--color-text-primary);
  }

  .admin-login button {
    padding: 0.625rem;
    background: var(--color-accent-primary);
    color: #fff;
    border: none;
    border-radius: 0.5rem;
    font-weight: 600;
    cursor: pointer;
  }

  .error {
    color: #ef4444;
    margin-top: 0.75rem;
    font-size: 0.875rem;
  }

  .admin {
    max-width: 960px;
    margin: 2rem auto;
    padding: 0 1rem;
  }

  .admin h1 {
    font-size: 1.5rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }

  .status-msg {
    background: #dcfce7;
    color: #166534;
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
    cursor: pointer;
    font-size: 0.875rem;
  }

  /* Legend */
  .legend {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
    font-size: 0.8125rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.375rem;
  }

  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
  }

  .dot.pending { background: #f59e0b; }
  .dot.confirmed { background: #22c55e; }
  .dot.rejected { background: #ef4444; }
  .dot.cancelled { background: #9ca3af; }

  /* Calendar */
  .cal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
  }

  .cal-header button {
    padding: 0.5rem 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-primary);
    border-radius: 0.375rem;
    cursor: pointer;
    color: var(--color-text-primary);
    font-size: 1rem;
  }

  .cal-header h2 {
    font-size: 1.25rem;
    font-weight: 700;
  }

  .loading {
    text-align: center;
    color: var(--color-text-muted);
    padding: 2rem;
  }

  .no-bookings {
    text-align: center;
    color: var(--color-text-muted);
    padding: 2rem;
    font-size: 0.875rem;
  }

  .admin-cal {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 1px;
    background: var(--color-border);
    border: 1px solid var(--color-border);
    border-radius: 0.5rem;
    overflow: hidden;
  }

  .cal-head {
    padding: 0.5rem;
    text-align: center;
    font-size: 0.75rem;
    font-weight: 600;
    background: var(--color-bg-secondary);
    color: var(--color-text-muted);
  }

  .cal-cell {
    min-height: 80px;
    padding: 0.25rem;
    background: var(--color-bg-primary);
    font-size: 0.75rem;
  }

  .cal-cell.empty {
    background: var(--color-bg-secondary);
  }

  .cell-day {
    font-weight: 600;
    font-size: 0.8125rem;
    margin-bottom: 0.25rem;
  }

  /* Booking chips */
  .booking-chip {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    padding: 0.125rem 0.25rem;
    border-radius: 0.25rem;
    font-size: 0.625rem;
    margin-bottom: 0.125rem;
    cursor: pointer;
    transition: opacity 0.15s;
  }

  .booking-chip:hover {
    opacity: 0.8;
  }

  .booking-chip.pending {
    background: #fef3c7;
    color: #92400e;
  }

  .booking-chip.confirmed {
    background: #dcfce7;
    color: #166534;
  }

  .booking-chip.rejected {
    background: #fee2e2;
    color: #991b1b;
    text-decoration: line-through;
  }

  .booking-chip.cancelled {
    background: #f3f4f6;
    color: #9ca3af;
    text-decoration: line-through;
  }

  .chip-time {
    font-weight: 600;
  }

  .chip-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  /* Modal */
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
    padding: 1rem;
  }

  .modal {
    background: var(--color-bg-primary);
    border-radius: 0.75rem;
    padding: 1.5rem;
    max-width: 480px;
    width: 100%;
    max-height: 80vh;
    overflow-y: auto;
  }

  .modal h3 {
    font-size: 1.125rem;
    font-weight: 700;
    margin-bottom: 1rem;
  }

  .modal p {
    font-size: 0.875rem;
    margin-bottom: 0.5rem;
  }

  .status-label {
    font-weight: 600;
    padding: 0.125rem 0.5rem;
    border-radius: 0.25rem;
    font-size: 0.8125rem;
  }

  .status-label.pending {
    background: #fef3c7;
    color: #92400e;
  }

  .status-label.confirmed {
    background: #dcfce7;
    color: #166534;
  }

  .status-label.rejected {
    background: #fee2e2;
    color: #991b1b;
  }

  .status-label.cancelled {
    background: #f3f4f6;
    color: #9ca3af;
  }

  .modal-actions {
    display: flex;
    gap: 0.75rem;
    margin-top: 1.25rem;
  }

  .cancel-btn {
    padding: 0.5rem 1rem;
    background: #ef4444;
    color: #fff;
    border: none;
    border-radius: 0.375rem;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.875rem;
  }

  .close-btn {
    padding: 0.5rem 1rem;
    border: 1px solid var(--color-border);
    background: var(--color-bg-primary);
    border-radius: 0.375rem;
    cursor: pointer;
    color: var(--color-text-primary);
    font-size: 0.875rem;
  }

  @media (max-width: 640px) {
    .admin-cal {
      font-size: 0.625rem;
    }

    .cal-cell {
      min-height: 60px;
    }
  }
</style>
