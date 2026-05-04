/** Reloj de referencia para el agente (Colombia), alineado con locale es-CO. */
export const BOGOTA_TZ = 'America/Bogota';

export function nowBogotaTimeHHMMSS(): string {
  return new Date().toLocaleTimeString('es-CO', {
    timeZone: BOGOTA_TZ,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

/** Fecha civil hoy en Bogotá, ideal para flight_search (YYYY-MM-DD). */
export function todayBogotaYYYYMMDD(): string {
  return new Date().toLocaleDateString('en-CA', { timeZone: BOGOTA_TZ });
}
